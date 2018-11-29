package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/juju/errors"
	"github.com/unrolled/render"
)

type processHandler struct {
	c  *ConfigManager
	rd *render.Render
}

func newProcessHandler(c *ConfigManager, rd *render.Render) *processHandler {
	return &processHandler{
		c:  c,
		rd: rd,
	}
}

// AddConfig adds config
func (p *processHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	cfg := &Config{}
	err := readJSON(r.Body, cfg)
	if err != nil {
		p.rd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	p.c.SetCfg(cfg.Path, cfg.Value)
	p.rd.JSON(w, http.StatusOK, nil)
}

// RemoveConfig removes config
func (p *processHandler) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	cfg := &Config{}
	err := readJSON(r.Body, cfg)
	if err != nil {
		p.rd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	p.c.RemoveCfg(cfg.Path)
	p.rd.JSON(w, http.StatusOK, nil)
}

// ListConfig lists config
func (p *processHandler) ListConfig(w http.ResponseWriter, r *http.Request) {
	p.rd.JSON(w, http.StatusOK, p.c.ListCfg())
}

func readJSON(r io.ReadCloser, data interface{}) error {
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Trace(err)
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
