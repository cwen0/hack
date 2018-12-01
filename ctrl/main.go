package main

import (
	"flag"

	"github.com/ngaut/log"
)

var (
	addr   string
	pdAddr string
)

func init() {
	flag.StringVar(&addr, "addr", ":10008", "operation manager address")
	flag.StringVar(&pdAddr, "pd-addr", "", "pd address")
}

func main() {
	operationManager := NewManager(addr, pdAddr)
	err := operationManager.Run()
	if err != nil {
		log.Fatalf("run operation manager failed %v", err)
	}
}
