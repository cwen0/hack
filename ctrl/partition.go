package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/juju/errors"
	"github.com/zhouqiang-cl/hack/network"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type networkCtl struct {
	toplogic *types.Topological
}

func newNetworkCtl(toplogic *types.Topological) *networkCtl {
	return &networkCtl{toplogic: toplogic}
}

func (n *networkCtl) start(kind types.PartitionKind) error {
	partition := &types.Partition{
		Kind: kind,
	}
	configs, err := network.GetProxyPartitionConfig(n.toplogic, partition)
	if err != nil {
		return errors.Trace(err)
	}
	// first empty
	for host := range configs {
		err := emptyPartition(host)
		if err != nil {
			return errors.Trace(err)
		}
	}

	// then set rule
	for host, cfg := range configs {
		err := doNetworkPartition(host, cfg)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func doNetworkPartition(host string, cfg *types.NetworkConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return errors.Trace(err)
	}
	url := fmt.Sprintf("http://%s:10008/config/network/partition/add", host)
	_, err = utils.DoPost(url, data)
	return errors.Trace(err)
}

func emptyPartition(host string) error {
	url := fmt.Sprintf("http://%s:10008/config/network/partition/remove", host)
	_, err := utils.DoPost(url, []byte{})
	return errors.Trace(err)
}
