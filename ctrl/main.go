package main

import (
	"flag"

	"github.com/ngaut/log"
)

var (
	addr     string
	pdAddr   string
	tidbAddr string
)

func init() {
	flag.StringVar(&addr, "addr", "172.16.30.12:10009", "operation manager address")
	flag.StringVar(&pdAddr, "pd-addr", "10.128.19.63:2379", "pd address")
	flag.StringVar(&tidbAddr, "tidb-addr", "10.128.26.206", "pd address")
}

func main() {
	flag.Parse()

	operationManager := NewManager(addr, pdAddr)
	err := operationManager.Run()
	if err != nil {
		log.Fatalf("run operation manager failed %v", err)
	}
}

func getAllPath() []string {
	return []string{
		"/tikvpb.Tikv/Coprocessor",
		"/tikvpb.Tikv/KvBatchGet",
		"/tikvpb.Tikv/KvBatchRollback",
		"/tikvpb.Tikv/KvCleanup",
		"/tikvpb.Tikv/KvCommit",
		"/tikvpb.Tikv/KvGet",
		"/tikvpb.Tikv/KvGC",
		"/tikvpb.Tikv/KvPrewrite",
		"/tikvpb.Tikv/KvResolveLock",
		"/tikvpb.Tikv/Raft",
		"/tikvpb.Tikv/SplitRegion",
		"/tikvpb.Tikv/KvScanLock",
	}
}
