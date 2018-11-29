package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	"github.com/ngaut/log"
)

var defaultPushMetricsInterval = 15 * time.Second

var (
	dbName      = flag.String("db", "test", "database name")
	addr      = flag.String("addr", "127.0.0.1:4000", "database name")
	accounts    = flag.Int("accounts", 1000000, "the number of accounts")
	interval    = flag.Duration("interval", 2*time.Second, "the interval")
	tables      = flag.Int("tables", 1, "the number of the tables")
	concurrency = flag.Int("concurrency", 200, "concurrency worker count")
	retryLimit  = flag.Int("retry-limit", 200, "retry count")
)

var (
	defaultVerifyTimeout = 6 * time.Hour
	remark               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXVZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXVZlkjsanksqiszndqpijdslnnq"
)

func main() {
	log.Info("[bank] start bank")
	flag.Parse()

	dbDSN := fmt.Sprintf("root:@tcp(%s)/%s", *addr, *dbName)
	db, err := OpenDB(dbDSN, *concurrency)
	if err != nil {
		log.Fatalf("[bank] create db client error %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-sc
		log.Infof("[bank] Got signal [%s] to exist.", sig)
		cancel()
		os.Exit(0)
	}()
	

	cfg := Config{
		NumAccounts: *accounts,
		Interval:    *interval,
		TableNum:    *tables,
		Concurrency: *concurrency,
	}
	bank := NewBankCase(&cfg)
	if err := bank.Initialize(ctx, db); err != nil {
		log.Fatalf("[bank] initial failed %v", err)
	}

	if err := bank.Execute(ctx, db); err != nil {
		log.Fatalf("[bank] return with error %v", err)
	}
}
