package main

import (
	"time"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ngaut/log"
	"github.com/juju/errors"
	"github.com/unrolled/render"
	"github.com/zhouqiang-cl/hack/types"
	"github.com/zhouqiang-cl/hack/utils"
)

var (
	storesPrefix     = "pd/api/v1/stores"
	storePrefix      = "pd/api/v1/store"
	schedulersPrefix = "pd/api/v1/schedulers"
)

type evictLeaderHandler struct {
	c  *Manager
	rd *render.Render
}

func newEvictLeaderHandler(c *Manager, rd *render.Render) *evictLeaderHandler {
	return &evictLeaderHandler{
		c:  c,
		rd: rd,
	}
}

func (e *evictLeaderHandler) EvictLeader(w http.ResponseWriter, r *http.Request) {
	tikvIP := r.URL.Query()["ip"]
	if len(tikvIP) == 0 {
		e.rd.JSON(w, http.StatusBadRequest, "miss parameter ip")
		return
	}

	err := doEvictLeader(tikvIP[0], e.c.pdAddr)
	if err != nil {
		e.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	e.rd.JSON(w, http.StatusOK, nil)
}

func doEvictLeader(tikvIP, pdAddr string) error {
	storesInfo, err := getStores(pdAddr)
	if err != nil {
		return errors.Trace(err)
	}

	var storeID uint64
	for _, store := range storesInfo.Stores {
		storeIP,ok  := utils.Resolve(store.Store.Address)
		if !ok {
			return errors.Errorf("address %s can not convert to ip", store.Store.Address)
		}
		if storeIP == tikvIP {
			storeID = store.Store.Id
		}
	}

	leaderEvictInfo := getLeaderEvictSchedulerInfo(storeID)
	apiURL := fmt.Sprintf("http://%s/%s", pdAddr, schedulersPrefix)
	data, err := json.Marshal(leaderEvictInfo)
	if err != nil {
		return err
	}

	_, err = utils.DoPost(apiURL, data)
	if err != nil {
		return err
	}

	for {
		storeInfo, err := getStore(storeID, pdAddr)
		if err != nil {
			return err
		}
		time.Sleep(3*time.Second)
		log.Debugf("leader count %d", storeInfo.Status.LeaderCount)
		if storeInfo.Status.LeaderCount == 0 {
			break
		}
	}

	return nil
}

func getStores(pdAddr string) (types.StoresInfo, error) {
	apiURL := fmt.Sprintf("http://%s/%s", pdAddr, storesPrefix)
	body, err := utils.DoGet(apiURL)
	if err != nil {
		return types.StoresInfo{}, errors.Trace(err)
	}

	storesInfo := types.StoresInfo{}
	err = json.Unmarshal(body, &storesInfo)
	if err != nil {
		return types.StoresInfo{}, errors.Trace(err)
	}

	return storesInfo, nil
}

func getStore(storeID uint64, pdAddr string) (types.StoreInfo, error) {
	apiURL := fmt.Sprintf("http://%s/%s/%d", pdAddr, storePrefix, storeID)
	body, err := utils.DoGet(apiURL)
	if err != nil {
		return types.StoreInfo{}, errors.Trace(err)
	}

	storeInfo := types.StoreInfo{}
	err = json.Unmarshal(body, &storeInfo)
	if err != nil {
		return types.StoreInfo{}, errors.Trace(err)
	}

	return storeInfo, nil
}

func getLeaderEvictSchedulerInfo(storeID uint64) *types.SchedulerInfo {
	return &types.SchedulerInfo{Name: "evict-leader-scheduler", StoreID: storeID}
}
