package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/unrolled/render"
	"github.com/zhouqiang-cl/hack/network"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
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
	kind := r.URL.Query()["kind"]
	if len(kind) == 0 {
		p.rd.JSON(w, http.StatusBadRequest, "miss parameter ip")
		return
	}
	localPartition := types.Partition{
		Kind: types.PartitionKind(kind[0]),
	}
	topology, err := getTopologyInfo(p.c.pdAddr)
	if err != nil {
		p.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	configs, err := network.GetProxyPartitionConfig(&topology, &localPartition)
	for name, cfg := range configs {
		log.Debugf("%s config is %+v", name, cfg)
	}
	if err != nil {
		p.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Infof("partition info: %+v", localPartition)
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
		partition: localPartition,
	}

	partition = localPartition

	logs.items = append(logs.items, Log{
		operation: OperationNetworkPartition,
		parameter: kind[0],
		timeStamp: time.Now().Unix(),
	})

	p.rd.JSON(w, http.StatusOK, nil)
}

func (p *partitionHandler) GetNetworkPartiton(w http.ResponseWriter, r *http.Request) {
	p.rd.JSON(w, http.StatusOK, partition)
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
