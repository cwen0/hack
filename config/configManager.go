package config

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/zhouqiang-cl/hack/types"
)

// APIPrefix is api prefix
const APIPrefix = "/config"

// Manager is the manager config.
type Manager struct {
	sync.RWMutex
	addr         string
	FailpointCfg map[string]string
	NetworkCfg   *types.NetworkConfig
	s            *http.Server
}

// NewManager creates the node with given address
func NewManager(addr string, cfg map[string]string) *Manager {
	n := &Manager{
		addr:         addr,
		FailpointCfg: cfg,
	}

	return n
}

// Run runs config manager
func (c *Manager) Run() error {
	c.s = &http.Server{
		Addr:    c.addr,
		Handler: c.createHandler(),
	}
	return c.s.ListenAndServe()
}

// Close closes the Node.
func (c *Manager) Close() {
	if c.s != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		c.s.Shutdown(ctx)
		cancel()
	}
}

// GetFailpointCfg gets failpoint cfg
func (c *Manager) GetFailpointCfg(name string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	rule, ok := c.FailpointCfg[name]
	return rule, ok
}

// SetFailpointCfg sets cfg
func (c *Manager) SetFailpointCfg(name string, rule string) {
	c.Lock()
	defer c.Unlock()
	c.FailpointCfg[name] = rule
}

// RemoveFailpointCfg removes failpoint cfg
func (c *Manager) RemoveFailpointCfg(name string) {
	c.Lock()
	defer c.Unlock()
	delete(c.FailpointCfg, name)
}

// ListFailpointCfg lists failpoint cfg
func (c *Manager) ListFailpointCfg() map[string]string {
	c.RLock()
	defer c.RUnlock()
	return c.FailpointCfg
}

// CleanFailpointCfg cleans failpoint cfg
func (c *Manager) CleanFailpointCfg() {
	c.Lock()
	defer c.Unlock()
	c.FailpointCfg = nil
}

// SetPartitionCfg sets network config
func (c *Manager) SetPartitionCfg(cfg *types.NetworkConfig) {
	c.Lock()
	defer c.Unlock()
	c.NetworkCfg = cfg
}

// GetPartitionCfg gets network config
func (c *Manager) GetPartitionCfg() *types.NetworkConfig {
	c.RLock()
	defer c.RUnlock()
	return c.NetworkCfg
}

// RemovePartitionCfg removes partition config
func (c *Manager) RemovePartitionCfg() {
	c.Lock()
	defer c.Unlock()
	c.NetworkCfg = nil
}

func (c *Manager) createHandler() http.Handler {
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

func (c *Manager) createRouter() *mux.Router {
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
