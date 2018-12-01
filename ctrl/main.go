package main

import (
	"flag"

	"github.com/zhouqiang-cl/hack/types"
	"github.com/ngaut/log"
)

var (
	addr   string
	pdAddr string
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:10009", "operation manager address")
	flag.StringVar(&pdAddr, "pd-addr", "", "pd address")
}

func main() {
	flag.Parse()

	operationManager := NewManager(addr, pdAddr)
	err := operationManager.Run()
	if err != nil {
		log.Fatalf("run operation manager failed %v", err)
	}
}

func getToplogic() *types.Topological {
	return &types.Topological{
		PD:   []string{"10.128.31.5"},
		TiDB: []string{"10.128.29.197"},
		TiKV: []string{"10.128.31.56", "10.128.31.62", "10.128.31.51"},
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
