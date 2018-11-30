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

// AddFailpointConfig adds failpoint config
func (p *processHandler) AddFailpointConfig(w http.ResponseWriter, r *http.Request) {
	cfg := &FailpointConfig{}
	err := readJSON(r.Body, cfg)
	if err != nil {
		p.rd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	p.c.SetFailpointCfg(cfg.Path, cfg.Value)
	p.rd.JSON(w, http.StatusOK, nil)
}

// RemoveFailpointConfig removes failpoint config
func (p *processHandler) RemoveFailpointConfig(w http.ResponseWriter, r *http.Request) {
	cfg := &FailpointConfig{}
	err := readJSON(r.Body, cfg)
	if err != nil {
		p.rd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	p.c.RemoveFailpointCfg(cfg.Path)
	p.rd.JSON(w, http.StatusOK, nil)
}

// ListFailpointConfig lists failpoint config
func (p *processHandler) ListFailpointConfig(w http.ResponseWriter, r *http.Request) {
	p.rd.JSON(w, http.StatusOK, p.c.ListFailpointCfg())
}

// CleanFailpointConfig cleans config
func (p *processHandler) CleanFailpointConfig(w http.ResponseWriter, r *http.Request) {
	p.c.CleanFailpointCfg()
	p.rd.JSON(w, http.StatusOK, nil)
}

// AddPartition adds partition
func (p *processHandler) AddPartitionConfig(w http.ResponseWriter, r *http.Request) {
	cfg := &NetworkConfig{}
	err := readJSON(r.Body, cfg)
	if err != nil {
		p.rd.JSON(w, http.StatusBadRequest, err.Error())
		return
	}
	p.c.SetPartitionCfg(cfg)
	p.rd.JSON(w, http.StatusOK, nil)
}

// RemovePartitionConfig removes config
func (p *processHandler) RemovePartitionConfig(w http.ResponseWriter, r *http.Request) {
	p.c.RemovePartitionCfg()
	p.rd.JSON(w, http.StatusOK, nil)
}

// GetPartitionConfig gets partition
func (p *processHandler) GetPartitionConfig(w http.ResponseWriter, r *http.Request) {
	cfg, ok := p.c.GetPartitionCfg()
	if ok {
		p.rd.JSON(w, http.StatusOK, cfg)
		return
	}

	p.rd.JSON(w, http.StatusOK, nil)
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
