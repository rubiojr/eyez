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

	proxy, err := eyez.New(ctx, &eyez.ProxyOptions{Port: 1080, CACert: *crt, CAKey: *key})
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
