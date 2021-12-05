package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rubiojr/eyez"
)

func main() {
	crt := flag.String("cacert", "", "CA Certificate")
	key := flag.String("cakey", "", "CA Key")
	port := flag.Int("port", 1080, "Proxy listening port")
	dbPath := flag.String("db", "records.db", "Database file path")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range signals {
			cancel()
		}
	}()

	proxy, err := eyez.New(ctx, &eyez.ProxyOptions{Port: *port, CACert: *crt, CAKey: *key, DBPath: *dbPath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating proxy: %s", err)
		os.Exit(1)
	}

	err = proxy.Serve()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error serving proxy: %s", err)
		os.Exit(1)
	}
}
