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
	flag.StringVar(&cmd, "cmd", "failpoint", "which command will use. now support failpoint/partition/record/replay, default is failpoint")
	flag.StringVar(&param, "param", "", "the param command will use")
}

// Ctrl is ctrl
type Ctrl struct {
	fpCtrl *failpointCtl
}

func newCtrl(toplogic *types.Topological) *Ctrl {
	return &Ctrl{
		fpCtrl: newFailpointCtl(toplogic),
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
	}
}

func getToplogic() *types.Topological {
	return &types.Topological{
		PD:   []string{"127.0.0.1"},
		TiDB: []string{"127.0.0.1", "127.0.0.2"},
		TiKV: []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"},
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
