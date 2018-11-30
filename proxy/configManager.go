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
	addr         string
	FailpointCfg map[string]string
	NetworkCfg   *NetworkConfig
	s            *http.Server
}

// NetworkConfig is network config
type NetworkConfig struct {
	Ingress  []string `json:"ingress"`
	Egress []string `json:"egress"`
}

// FailpointConfig is failpoint config
type FailpointConfig struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

// NewConfigManager creates the node with given address
func NewConfigManager(addr string, cfg map[string]string) *ConfigManager {
	n := &ConfigManager{
		addr:         addr,
		FailpointCfg: cfg,
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

// GetFailpointCfg gets failpoint cfg
func (c *ConfigManager) GetFailpointCfg(name string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	rule, ok := c.FailpointCfg[name]
	return rule, ok
}

// SetFailpointCfg sets cfg
func (c *ConfigManager) SetFailpointCfg(name string, rule string) {
	c.Lock()
	defer c.Unlock()
	c.FailpointCfg[name] = rule
}

// RemoveFailpointCfg removes failpoint cfg
func (c *ConfigManager) RemoveFailpointCfg(name string) {
	c.Lock()
	defer c.Unlock()
	delete(c.FailpointCfg, name)
}

// ListFailpointCfg lists failpoint cfg
func (c *ConfigManager) ListFailpointCfg() map[string]string {
	c.RLock()
	defer c.RUnlock()
	return c.FailpointCfg
}

// CleanFailpointCfg cleans failpoint cfg
func (c *ConfigManager) CleanFailpointCfg() {
	c.Lock()
	defer c.Unlock()
	c.FailpointCfg = nil
}

// SetPartitionCfg sets network config
func (c *ConfigManager) SetPartitionCfg(cfg *NetworkConfig) {
	c.Lock()
	defer c.Unlock()
	c.NetworkCfg = cfg
}

// GetPartitionCfg gets network config
func (c *ConfigManager) GetPartitionCfg() *NetworkConfig {
	c.RLock()
	defer c.RUnlock()
	return c.NetworkCfg
}

// RemovePartitionCfg removes partition config
func (c *ConfigManager) RemovePartitionCfg() {
	c.Lock()
	defer c.Unlock()
	c.NetworkCfg = nil
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

	// failpoint
	router.HandleFunc("/failpoint/add", processHandler.AddFailpointConfig).Methods("POST")
	router.HandleFunc("/failpoint/remove", processHandler.RemoveFailpointConfig).Methods("POST")
	router.HandleFunc("/failpoint/clean", processHandler.CleanFailpointConfig).Methods("POST")
	router.HandleFunc("/failpoint", processHandler.ListFailpointConfig).Methods("GET")

	// network partition
	router.HandleFunc("/network/partition/add", processHandler.AddPartitionConfig).Methods("POST")
	router.HandleFunc("/network/partition/remove", processHandler.RemovePartitionConfig).Methods("POST")
	router.HandleFunc("/network/partition", processHandler.GetPartitionConfig).Methods("GET")

	return router
}
