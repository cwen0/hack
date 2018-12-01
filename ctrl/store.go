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
	storesInfo, err := getStores(pdAddr)
	if err != nil {
		s.rd.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.rd.JSON(w, http.StatusOK, storesInfo)
}
