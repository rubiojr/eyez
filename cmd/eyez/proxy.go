package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rubiojr/eyez"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range signals {
			cancel()
		}
	}()

	proxy, err := eyez.New(ctx, &eyez.ProxyOptions{Port: 1080})
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
