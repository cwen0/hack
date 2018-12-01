package main

import (
	"github.com/unrolled/render"
	"net/http"
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
	} else {
		leaderCnt := 0
		for _, store := range storesInfo.Stores {
			if store.Store.Address == tikvIP[0] {
				leaderCnt = store.Status.LeaderCount
			}
		}
		s.rd.JSON(w, http.StatusOK, leaderCnt)
		return
	}


}

func (s *storeHandler) GetStoreLeader(w http.ResponseWriter, r *http.Request) {

}
