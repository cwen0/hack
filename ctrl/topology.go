package main

import (
	"encoding/json"
	"fmt"
	"github.com/unrolled/render"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
	"net/http"
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
	doTopology(t.c.pdAddr)
}

func getMembers(pdAddr string) (*types.MembersInfo, error) {
	apiURL := fmt.Sprintf("%s/%s", pdAddr, membersPrefix)
	body, err := utils.DoGet(apiURL)
	if err != nil {
		return nil, err
	}

	membersInfo := types.MembersInfo{}
	err = json.Unmarshal(body, &membersInfo)
	if err != nil {
		return nil, err
	}

	return &membersInfo, nil
}

func doTopology(pdAddr string) (*types.Topological, error) {
	var topoInfo *types.Topological
	topoInfo, err := getTopologyInfo(pdAddr)
	if err != nil {
		return nil, err
	}
	return topoInfo, nil
}


func getTopologyInfo(pdAddr string) (*types.Topological, error){
	var topologyInfo types.Topological
	storesInfo, err := getStores(pdAddr)
	if err != nil {
		return nil, err
	}
	membersInfo, err := getMembers(pdAddr)

	for _, store := range storesInfo.Stores {
		topologyInfo.TiKV = append(topologyInfo.TiKV, store.Store.Address)
	}
	for _, member := range membersInfo.Members{
		topologyInfo.PD = append(topologyInfo.PD, member.Name)
	}

	return &topologyInfo, nil
}
