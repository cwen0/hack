package main

import (
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
	"github.com/juju/errors"
	"fmt"
	"net/http"
	"time"
	"encoding/json"
)

type topology struct {
	url        string
	httpClient *http.Client
}


func (e *topology) start() error {
	return errors.Trace(e.doTopology())
}


func newTopplogyCtl(url string, timeout time.Duration) *topology {
	return &topology{
		url:        url,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (e *topology) getStores() (*types.StoresInfo,error) {
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


func (e *topology) getMembers() (*types.MembersInfo,error) {
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

func (e *topology) doTopology() (error) {
	for {
		_, err := e.getTopologyInfo()
		if err != nil {
			return err
		}
	}

	return nil
}


func (e *topology)getTopologyInfo() (*types.TopologyInfo, error){
	storesInfo, _ := e.getStores()
	membersInfo, _ := e.getMembers()
	topoInfo := types.TopologyInfo{StoresInfo: storesInfo, MembersInfo:membersInfo}
	return &topoInfo, nil
}