package main

import (
	"context"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/zhouqiang-cl/hack/config"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	// "google.golang.org/grpc/connectivity"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

// ProxyHandler is proxy handler
type ProxyHandler struct {
	ctx          context.Context
	cfgManager   *config.Manager
	upstream     string
	upstreamConn *grpc.ClientConn
}

// NewProxyHandler creates new proxy handler
func NewProxyHandler(ctx context.Context, upstream string, cfgManager *config.Manager) (*ProxyHandler, error) {
	streamer := &ProxyHandler{
		ctx:        ctx,
		cfgManager: cfgManager,
		upstream:   upstream,
	}

	var err error

	streamer.upstreamConn, err = grpc.DialContext(ctx, upstream, grpc.WithInsecure(), grpc.WithCodec(Codec()))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return streamer, nil
}

// StreamHandler returns stream handler
func (p *ProxyHandler) StreamHandler() grpc.StreamHandler {
	return p.handler
}

func (p *ProxyHandler) handler(srv interface{}, serverStream grpc.ServerStream) error {
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}
	log.Infof("full name %s", fullMethodName)
	clientStream, err := grpc.NewClientStream(p.ctx, clientStreamDescForProxying, p.upstreamConn, fullMethodName)
	if err != nil {
		log.Errorf("failed to create client %+v, state is %v", err, p.upstreamConn.GetState())
		return err
	}

	pe, ok := peer.FromContext(serverStream.Context())
	if !ok {
		log.Error("get peer failed")
	}
	log.Debugf("perr addr: %s", pe.Addr.String())

	s2cErrChan := p.forwardServerToClient(serverStream, clientStream)
	c2sErrChan := p.forwardClientToServer(clientStream, serverStream)

	// We don't know which side is going to stop sending first, so we need a select between the two.
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan:
			if s2cErr == io.EOF {
				// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore./
				// the clientStream>serverStream may continue pumping though.
				clientStream.CloseSend()
				break
			} else {
				// however, we may have gotten a receive error (stream disconnected, a read error etc) in which case we need
				// to cancel the clientStream to the backend, let all of its goroutines be freed up by the CancelFunc and
				// exit with an error to the stack
				return grpc.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// This happens when the clientStream has nothing else to offer (io.EOF), returned a gRPC error. In those two
			// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
			// will be nil.
			serverStream.SetTrailer(clientStream.Trailer())
			// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}
	}
	return grpc.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
}

func (p *ProxyHandler) forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		// f := &frame{}
		for i := 0; ; i++ {
			if err := p.handleOutRequest(i, src, dst); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (p *ProxyHandler) forwardServerToClient(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		for i := 0; ; i++ {
			if err := p.handleInRequest(src, dst); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

// handleRequest try to apply config
func (p *ProxyHandler) handleInRequest(src grpc.ServerStream, dst grpc.ClientStream) error {
	methodName, ok := grpc.MethodFromServerStream(src)
	if !ok {
		return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}

	cfg, ok := p.cfgManager.GetPartitionCfg()
	if ok && len(cfg.Ingress) > 0 {
		log.Infof("start to filter ingress request")
		log.Infof("ingress allow hosts: %v", cfg.Ingress)
		if err := p.processIngressNetwork(src, dst, cfg); err != nil {
			return err
		}
		return nil
	}

	rule, ok := p.cfgManager.GetFailpointCfg(methodName)
	if !ok {
		return p.processNormal(src, dst)
	}

	return p.processWithRule(src, dst, rule)
}

func (p *ProxyHandler) handleOutRequest(index int, src grpc.ClientStream, dst grpc.ServerStream) error {
	cfg, ok := p.cfgManager.GetPartitionCfg()
	if ok && len(cfg.Egress) > 0 {
		log.Infof("start to filter egress request")
		log.Infof("egress allow hosts: %v", cfg.Egress)
		if err := p.processEgressNetwork(index, src, dst, cfg); err != nil {
			return err
		}
		return nil
	}

	return p.processOutNormal(index, src, dst)
}

func (p *ProxyHandler) processNormal(src grpc.ServerStream, dst grpc.ClientStream) error {
	f := &frame{}
	err := src.RecvMsg(f)
	if err != nil {
		// can not use error.Trace for eof
		return err
	}
	//log.Debugf("data is %+v", f)

	err = dst.SendMsg(f)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyHandler) processOutNormal(index int, src grpc.ClientStream, dst grpc.ServerStream) error {
	f := &frame{}
	if err := src.RecvMsg(f); err != nil {
		return err
	}
	if index == 0 {
		// This is a bit of a hack, but client to server headers are only readable after first client msg is
		// received but must be written to server stream before the first msg is flushed.
		// This is the only place to do it nicely.
		md, err := src.Header()
		if err != nil {
			return err
		}
		if err := dst.SendHeader(md); err != nil {
			return err
		}
	}
	if err := dst.SendMsg(f); err != nil {
		return err
	}

	return nil
}

func (p *ProxyHandler) processWithRule(src grpc.ServerStream, dst grpc.ClientStream, ruleStr string) error {
	f := &frame{}
	err := src.RecvMsg(f)
	if err != nil {
		// can not use error.Trace for eof
		return err
	}

	rules := getRulesFromRuleStr(ruleStr)
	for _, rule := range rules {
		n := rand.Intn(100)
		if n >= rule.Pct {
			continue
		}
		switch rule.Action {
		case "delay":
			millisecond, err := strconv.ParseInt(rule.ActionArgs, 10, 64)
			if err != nil {
				return errors.Trace(err)
			}
			log.Infof("sleep %d ms", millisecond)
			time.Sleep(time.Duration(millisecond) * time.Millisecond)
		case "timeout":
			log.Infof("sleep 10 minutes for timeout")
			time.Sleep(10 * time.Minute)
		}
	}

	err = dst.SendMsg(f)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyHandler) processIngressNetwork(src grpc.ServerStream, dst grpc.ClientStream, cfg *types.NetworkConfig) error {
	pe, ok := peer.FromContext(src.Context())
	if !ok {
		log.Error("get peer failed")
	}

	log.Infof("Ingress perr addr: %s", pe.Addr.String())

	ingressIP, err := utils.GetIP(pe.Addr.String())
	if err != nil {
		return err
	}

	if utils.MatchInArray(cfg.Egress, ingressIP) {
		return p.processNormal(src, dst)
	}

	log.Infof("drop ingress request: %s", pe.Addr.String())

	return nil
}

func (p *ProxyHandler) processEgressNetwork(index int, src grpc.ClientStream, dst grpc.ServerStream, cfg *types.NetworkConfig) error {
	pe, ok := peer.FromContext(dst.Context())
	if !ok {
		log.Error("get peer failed")
	}

	log.Infof("Egress perr addr: %s", pe.Addr.String())
	egressIP, err := utils.GetIP(pe.Addr.String())
	if err != nil {
		return err
	}

	if utils.MatchInArray(cfg.Egress, egressIP) {
		return p.processOutNormal(index, src, dst)
	}

	log.Infof("drop egress request: %s", pe.Addr.String())

	return nil
}
