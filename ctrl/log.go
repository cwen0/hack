package main

import (
	"github.com/unrolled/render"
	"net/http"
)

type logHandler struct {
	c  *Manager
	rd *render.Render
}

func newLogHandler(c *Manager, rd *render.Render) *logHandler {
	return &logHandler{
		c:  c,
		rd: rd,
	}
}

func (f *logHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	f.rd.JSON(w, http.StatusOK, logs)
}
