package main

import (
	"fmt"
	"net/http"

	"github.com/zhouqiang-cl/hack/utils"
	"github.com/unrolled/render"
)

type storeHandler struct {
	c  *Manager
	rd *render.Render
}

func newStoreHandler(c *Manager, rd *render.Render) *storeHandler {
	return &storeHandler{
		c:  c,
		rd: rd,
	}
}

func (s *storeHandler) GetStores(w http.ResponseWriter, r *http.Request) {
	tikvIP := r.URL.Query()["ip"]
	storesInfo, err := getStores(pdAddr)
	if err != nil {
		s.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(tikvIP) == 0 {
		s.rd.JSON(w, http.StatusOK, storesInfo)
		return
	}
	leaderCnt := 0
	for _, store := range storesInfo.Stores {
		storeIP, ok := utils.Resolve(store.Store.Address)
		if !ok {
			s.rd.JSON(w, http.StatusInternalServerError, fmt.Sprintf("can not resolve %s", store.Store.Address))
			return
		}
		if storeIP == tikvIP[0] {
			leaderCnt = store.Status.LeaderCount
		}
	}
	s.rd.JSON(w, http.StatusOK, leaderCnt)
	return

}
