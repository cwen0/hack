package main

import (
	"flag"

	"github.com/ngaut/log"
	"github.com/zhouqiang-cl/hack/types"
)

var (
	cmd   string
	param string
)

func init() {
	flag.StringVar(&cmd, "cmd", "failpoint", "which command will use. now support failpoint/partition/record/replay/config/distributary, default is failpoint")
	flag.StringVar(&param, "param", "", "the param command will use")
}

// Ctrl is ctrl
type Ctrl struct {
	fpCtrl *failpointCtl
	npCtrl *networkCtl
}

func newCtrl(toplogic *types.Topological) *Ctrl {
	return &Ctrl{
		fpCtrl: newFailpointCtl(toplogic),
		npCtrl: newNetworkCtl(toplogic),
	}
}

func main() {
	flag.Parse()

	toplogic := getToplogic()
	ctrl := newCtrl(toplogic)
	switch cmd {
	case "failpoint":
		err := ctrl.fpCtrl.start(param)
		if err != nil {
			log.Fatal("failpoint failed %+v", err)
		}
	case "network":
		err := ctrl.npCtrl.start(types.PartitionKind(param))
		if err != nil {
			log.Fatalf("network failed %+v", err)
		}
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
		"/tikvpb.Tikv/KvCommit",
		"/tikvpb.Tikv/KvGet",
		"/tikvpb.Tikv/KvPrewrite",
		"/tikvpb.Tikv/KvScanLock",
	}
}
