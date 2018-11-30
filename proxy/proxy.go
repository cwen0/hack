package main

import (
	"context"
	"flag"
	"net"

	"github.com/ngaut/log"
	"github.com/zhouqiang-cl/hack/config"
	"google.golang.org/grpc"
)

var (
	upstream   string
	listenAddr string
	configAddr string
)

func init() {
	flag.StringVar(&upstream, "upstream", "127.0.0.1:20161", "upstream port")
	flag.StringVar(&listenAddr, "listen-addr", ":20160", "serve port")
	flag.StringVar(&configAddr, "config-listen-addr", ":10008", "config listen addr")
}

func main() {
	flag.Parse()

	cfg := defaultCfg()

	cfgManager := config.NewManager(configAddr, cfg)
	go func() {
		err := cfgManager.Run()
		if err != nil {
			log.Fatalf("run config failed %v", err)
		}
	}()
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx := context.Background()

	proxyHandler, err := NewProxyHandler(ctx, upstream, cfgManager)
	if err != nil {
		log.Fatalf("failed to setup proxy: %v", err)
	}
	s := grpc.NewServer(grpc.CustomCodec(Codec()),
		grpc.UnknownServiceHandler(proxyHandler.StreamHandler()))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func defaultCfg() map[string]string {
	cfg := make(map[string]string)
	cfg["/tikvpb.Tikv/Coprocessor"] = "pct(5)->delay(100)|pct(1)->timeout()"
	cfg["/tikvpb.Tikv/KvBatchGet"] = "pct(5)->delay(10)|pct(1)->timeout()"
	cfg["/tikvpb.Tikv/KvCommit"] = "pct(5)->delay(10)|pct(1)->timeout()"
	cfg["/tikvpb.Tikv/KvGet"] = "pct(5)->delay(10)|pct(1)->timeout()"
	cfg["/tikvpb.Tikv/KvPrewrite"] = "pct(5)->delay(10)|pct(1)->timeout()"
	cfg["/tikvpb.Tikv/KvScanLock"] = "pct(5)->delay(10)|pct(1)->timeout()"
	return cfg
}
