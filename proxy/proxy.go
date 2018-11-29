package main

import (
	"context"
	"flag"
	"log"
	"net"

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
	flag.StringVar(&configAddr, "config-listen-addr", ":100080", "config listen addr")
}

func main() {
	flag.Parse()

	cfg := defaultCfg()

	cfgManager := NewConfigManager(configAddr, cfg)
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
	cfg["/tikvpb.Tikv/Coprocessor"] = "rand(5)->delay(100)|rand(1)->timeout()"
	cfg["/tikvpb.Tikv/KvBatchGet"] = "rand(5)->delay(10)|rand(1)->timeout()"
	cfg["/tikvpb.Tikv/KvCommit"] = "rand(5)->delay(10)|rand(1)->timeout()"
	cfg["/tikvpb.Tikv/KvGet"] = "rand(5)->delay(10)|rand(1)->timeout()"
	cfg["/tikvpb.Tikv/KvPrewrite"] = "rand(5)->delay(10)|rand(1)->timeout()"
	cfg["/tikvpb.Tikv/KvScanLock"] = "rand(5)->delay(10)|rand(1)->timeout()"
	return cfg
}
