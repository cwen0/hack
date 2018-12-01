package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/juju/errors"
	"github.com/zhouqiang-cl/hack/network"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
	"github.com/ngaut/log"
	"github.com/unrolled/render"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type partitionHandler struct {
	c  *Manager
	rd *render.Render
}

func newPartitionHandler(c *Manager, rd *render.Render) *partitionHandler {
	return &partitionHandler{
		c:  c,
		rd: rd,
	}
}

func (p *partitionHandler) CreateNetworkPartition(w http.ResponseWriter, r *http.Request) {
	kind := r.URL.Query()["kind"][0]
	partition := &types.Partition{
		Kind: types.PartitionKind(kind),
	}
	configs, err := network.GetProxyPartitionConfig(getToplogic(), partition)
	for name, cfg := range configs {
		log.Debugf("%s config is %+v", name, cfg)
	}
	if err != nil {
		p.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Infof("partition info: %+v", partition)
	// first empty
	for host := range configs {
		err := emptyPartition(host)
		if err != nil {
			p.rd.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// then set rule
	for host, cfg := range configs {
		err := doNetworkPartition(host, cfg)
		if err != nil {
			p.rd.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	state = State{
		operation: StateNetworkPartition,
		partition: *partition,
	}


	p.rd.JSON(w, http.StatusOK, nil)
}

func (p *partitionHandler) GetNetworkPartiton(w http.ResponseWriter, r *http.Request) {

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
