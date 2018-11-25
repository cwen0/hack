package main

import (
	"log"
	"io"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"github.com/juju/errors"
)

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

// ProxyHandler is proxy handler
type ProxyHandler struct {
	ctx context.Context
	cfg map[string]string
	upstream string
	upstreamConn *grpc.ClientConn
}

// NewProxyHandler creates new proxy handler
func NewProxyHandler(ctx context.Context, cfg map[string]string, upstream string) (*ProxyHandler,error) {
	streamer := &ProxyHandler{
		ctx: ctx,
		cfg:cfg,
		upstream: upstream,
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

func (p *ProxyHandler)handler(srv interface{}, serverStream grpc.ServerStream) error {
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return grpc.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}
	log.Printf("full name %s", fullMethodName)
	clientStream, err := grpc.NewClientStream(p.ctx, clientStreamDescForProxying, p.upstreamConn, fullMethodName)
	if err != nil {
		return err
	}

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
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if i == 0 {
				// This is a bit of a hack, but client to server headers are only readable after first client msg is
				// received but must be written to server stream before the first msg is flushed.
				// This is the only place to do it nicely.
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
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
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}