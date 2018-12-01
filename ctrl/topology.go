package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/juju/errors"
	"github.com/unrolled/render"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
)

var membersPrefix = "pd/api/v1/members"

type topologyHandler struct {
	c  *Manager
	rd *render.Render
}

func newTopologynHandler(c *Manager, rd *render.Render) *topologyHandler {
	return &topologyHandler{
		c:  c,
		rd: rd,
	}
}

func (t *topologyHandler) GetTopology(w http.ResponseWriter, r *http.Request) {
	topology, err := getTopologyInfo(t.c.pdAddr)
	if err != nil {
		t.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	t.rd.JSON(w, http.StatusOK, topology)
}

func getMembers(pdAddr string) (types.MembersInfo, error) {
	apiURL := fmt.Sprintf("http://%s/%s", pdAddr, membersPrefix)
	body, err := utils.DoGet(apiURL)
	if err != nil {
		return types.MembersInfo{}, err
	}

	membersInfo := types.MembersInfo{}
	err = json.Unmarshal(body, &membersInfo)
	if err != nil {
		return types.MembersInfo{}, err
	}

	return membersInfo, nil
}

func getTopologyInfo(pdAddr string) (types.Topological, error) {
	var topologyInfo types.Topological
	storesInfo, err := getStores(pdAddr)
	if err != nil {
		return types.Topological{}, errors.Trace(err)
	}
	membersInfo, err := getMembers(pdAddr)
	if err != nil {
		return types.Topological{}, errors.Trace(err)
	}

	for _, store := range storesInfo.Stores {
		topologyInfo.TiKV = append(topologyInfo.TiKV, store.Store.Address)
	}
	for _, member := range membersInfo.Members {
		topologyInfo.PD = append(topologyInfo.PD, member.Name)
	}

	return topologyInfo, nil
}
