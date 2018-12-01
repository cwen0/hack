package cmd

import (
	"flag"
	"fmt"
	"github.com/ngaut/log"
	"github.com/zhouqiang-cl/hack/utils"
)

var (
	cmd         string
	param       string
	managerAddr = "127.0.0.1:10009"
)

func init() {
	flag.StringVar(&cmd, "cmd", "failpoint", "which command will use. now support failpoint/partition/record/replay/config/distributary, default is failpoint")
	flag.StringVar(&param, "param", "", "the param command will use")
}

func main() {
	flag.Parse()

	switch cmd {
	case "failpoint":
		url := fmt.Sprintf("http://%s/failpoint?type=%s", managerAddr, param)
		_, err := utils.DoPost(url, nil)
		if err != nil {
			log.Fatal("failpoint failed %+v", err)
		}
	case "network":
		url := fmt.Sprintf("http://%s/partition?kind=%s", managerAddr, param)
		_, err := utils.DoPost(url, nil)
		if err != nil {
			log.Fatalf("network failed %+v", err)
		}
	case "evict_leader":
		url := fmt.Sprintf("http://%s/evictleader?ip=%s", managerAddr, param)
		_, err := utils.DoPost(url, nil)
		if err != nil {
			log.Fatalf("evict leader failed %+v", err)
		}
	case "topology":
		url := fmt.Sprintf("http://%s/topology", managerAddr, param)
		_, err := utils.DoGet(url)
		if err != nil {
			log.Fatalf("get topology failed %+v", err)
		}
	}
}
