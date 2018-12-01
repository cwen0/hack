package main

import (
	"github.com/unrolled/render"
	"net/http"
)

type stateHandler struct {
	c  *Manager
	rd *render.Render
}

func newStateHandler(c *Manager, rd *render.Render) *stateHandler {
	return &stateHandler{
		c:  c,
		rd: rd,
	}
}

func (f *stateHandler) GetState(w http.ResponseWriter, r *http.Request) {
	f.rd.JSON(w, http.StatusOK, state)
}
