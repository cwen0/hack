package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/zhouqiang-cl/hack/utils"
	"github.com/zhouqiang-cl/hack/types"
)

var (
	cmd         string
	param       string
	managerAddr = "127.0.0.1:10009"
)

func init() {
	flag.StringVar(&cmd, "cmdline", "failpoint", "which command will use. now support failpoint/partition/record/replay/config/distributary, default is failpoint")
	flag.StringVar(&param, "param", "", "the param command will use")
}

func main() {
	flag.Parse()

	switch cmd {
	case "failpoint":
		url := fmt.Sprintf("http://%s/operation/failpoint?type=%s", managerAddr, param)
		_, err := utils.DoPost(url, []byte{})
		if err != nil {
			log.Fatal("failpoint failed %+v", errors.ErrorStack(err))
		}
	case "network":
		url := fmt.Sprintf("http://%s/operation/partition?kind=%s", managerAddr, param)
		_, err := utils.DoPost(url, []byte{})
		if err != nil {
			log.Fatalf("network failed %+v", errors.ErrorStack(err))
		}
	case "evict_leader":
		url := fmt.Sprintf("http://%s/operation/evictleader?ip=%s", managerAddr, param)
		_, err := utils.DoPost(url, []byte{})
		if err != nil {
			log.Fatalf("evict leader failed %+v", errors.ErrorStack(err))
		}
	case "topology":
		url := fmt.Sprintf("http://%s/operation/topology", managerAddr)
		data, err := utils.DoGet(url)
		if err != nil {
			log.Fatalf("get topology failed %+v", errors.ErrorStack(err))
		}
		top := &types.Topological{}
		err = json.Unmarshal(data, top)
		if err != nil {
			log.Fatalf("unmarshal topology failed %+v", errors.ErrorStack(err))
		}

		for _, ip := range top.PD{
			log.Infof("pd: %v", ip)
		}

		for _, ip := range top.TiKV{
			log.Infof("tikv: %v", ip)
		}

		for _, ip := range top.TiDB{
			log.Infof("tidb: %v", ip)
		}
	}
}
