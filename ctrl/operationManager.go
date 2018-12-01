package main

import (
	"context"
	"github.com/zhouqiang-cl/hack/types"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// prefix
const (
	APIPrefix                 = "/operation"
	StateNetworkPartition     = "network_partition"
	StateFailpoint            = "failpoint"
	OperationNetworkPartition = "network_partition"
	OperationFailpoint        = "failpoint"
)

// State is state
type State struct {
	operation string
	partition types.Partition
}

// Logs is logs
type Logs struct {
	Items []Log
}

// Log is log
type Log struct {
	operation string
	parameter string
	timeStamp int64
}

var state State
var logs Logs
var partition types.Partition

// Manager is the operation manager.
type Manager struct {
	sync.RWMutex
	addr   string
	pdAddr string
	s      *http.Server
}

// NewManager creates the node with given address
func NewManager(addr, pdAddr string) *Manager {
	n := &Manager{
		addr:   addr,
		pdAddr: pdAddr,
	}

	return n
}

// Run runs operation manager
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

	failpointHandler := newFailpointHandler(c, rd)
	partitionHandler := newPartitionHandler(c, rd)
	topologyHandler := newTopologynHandler(c, rd)
	evictLeaderHandler := newEvictLeaderHandler(c, rd)
	logHandler := newLogHandler(c, rd)
	stateHandler := newStateHandler(c, rd)
	storeHandler := newStoreHandler(c, rd)

	// failpoint route
	router.HandleFunc("/failpoint", failpointHandler.CreateFailpoint).Methods("POST")

	// network partition route
	router.HandleFunc("/partition", partitionHandler.CreateNetworkPartition).Methods("POST")
	router.HandleFunc("/partition/clean", partitionHandler.CleanNetworkPartition).Methods("POST")
	router.HandleFunc("/partition", partitionHandler.GetNetworkPartiton).Methods("GET")

	// topology route
	router.HandleFunc("/topology", topologyHandler.GetTopology).Methods("GET")

	// evict leader route
	router.HandleFunc("/evictleader", evictLeaderHandler.EvictLeader).Methods("POST")

	// log route
	router.HandleFunc("/log", logHandler.GetLogs).Methods("GET")

	// state route
	router.HandleFunc("/state", stateHandler.GetState).Methods("GET")

	//store route
	router.HandleFunc("/store", storeHandler.GetStores).Methods("GET")

	return router
}
