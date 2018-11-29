package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// APIPrefix is api prefix
const APIPrefix = "/config"

// ConfigManager is the manager config.
type ConfigManager struct {
	sync.RWMutex
	addr string
	Cfg  map[string]string
	s    *http.Server
}

// Config is config
type Config struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

// NewConfigManager creates the node with given address
func NewConfigManager(addr string, cfg map[string]string) *ConfigManager {
	n := &ConfigManager{
		addr: addr,
		Cfg:  cfg,
	}

	return n
}

// Run runs config manager
func (c *ConfigManager) Run() error {
	c.s = &http.Server{
		Addr:    c.addr,
		Handler: c.createHandler(),
	}
	return c.s.ListenAndServe()
}

// Close closes the Node.
func (c *ConfigManager) Close() {
	if c.s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		c.s.Shutdown(ctx)
		cancel()
	}
}

// GetCfg gets cfg
func (c *ConfigManager) GetCfg(name string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	rule, ok := c.Cfg[name]
	return rule, ok
}

// SetCfg sets cfg
func (c *ConfigManager) SetCfg(name string, rule string) {
	c.Lock()
	defer c.Unlock()
	c.Cfg[name] = rule
}

// RemoveCfg removes cfg
func (c *ConfigManager) RemoveCfg(name string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Cfg, name)
}

// ListCfg lists cfg
func (c *ConfigManager) ListCfg() map[string]string {
	c.RLock()
	defer c.RUnlock()
	return c.Cfg
}

func (c *ConfigManager) createHandler() http.Handler {
	engine := negroni.New()
	recover := negroni.NewRecovery()
	engine.Use(recover)

	router := mux.NewRouter()
	subRouter := c.createRouter()
	router.PathPrefix(APIPrefix).Handler(
		negroni.New(negroni.Wrap(subRouter)),
	)

	engine.UseHandler(router)
	return engine
}

func (c *ConfigManager) createRouter() *mux.Router {
	rd := render.New(render.Options{
		IndentJSON: true,
	})

	router := mux.NewRouter().PathPrefix(APIPrefix).Subrouter()

	processHandler := newProcessHandler(c, rd)
	router.HandleFunc("/add", processHandler.AddConfig).Methods("POST")
	router.HandleFunc("/remove", processHandler.RemoveConfig).Methods("POST")
	router.HandleFunc("/", processHandler.ListConfig).Methods("GET")

	return router
}
