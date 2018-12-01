package main

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errors"
	"github.com/unrolled/render"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
	"net/http"
	"time"
)

var membersPrefix = "pd/api/v1/members"

type topology struct {
	url        string
	httpClient *http.Client
}

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

func (f *topologyHandler) GetTopology(w http.ResponseWriter, r *http.Request) {

}

func (e *topology) start(tikvIP string) (*types.Topological, error) {
	topo, err := e.doTopology()
	return topo ,errors.Trace(err)
}

func newTopplogyCtl(url string, timeout time.Duration) *topology {
	return &topology{
		url:        url,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (e *topology) getStores() (*types.StoresInfo, error) {
	apiURL := fmt.Sprintf("%s/%s", e.url, storesPrefix)
	body, err := utils.DoGet(apiURL)
	if err != nil {
		return nil, err
	}

	storesInfo := types.StoresInfo{}
	err = json.Unmarshal(body, &storesInfo)
	if err != nil {
		return nil, err
	}

	return &storesInfo, nil
}

func (e *topology) getMembers() (*types.MembersInfo, error) {
	apiURL := fmt.Sprintf("%s/%s", e.url, membersPrefix)
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

func (e *topology) doTopology() (*types.Topological, error) {
	var topoInfo *types.Topological
	topoInfo, err := e.getTopologyInfo()
	if err != nil {
		return nil, err
	}
	return topoInfo, nil
}


func (e *topology)getTopologyInfo() (*types.Topological, error){
	var topologyInfo *types.Topological
	storesInfo, _ := e.getStores()
	membersInfo, _ := e.getMembers()

	for _, store := range storesInfo.Stores {
		topologyInfo.TiKV = append(topologyInfo.TiKV, store.Store.Address)
	}
	for _, member := range membersInfo.Members{
		topologyInfo.PD = append(topologyInfo.PD, member.Name)
	}
	return topologyInfo, nil
}
