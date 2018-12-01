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
	flag.StringVar(&addr, "addr", "127.0.0.1:10009", "operation manager address")
	flag.StringVar(&pdAddr, "pd-addr", "", "pd address")
	flag.StringVar(&tidbAddr, "tidb-addr", "10.128.29.197", "pd address")
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
