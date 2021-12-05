package main

import (
	"context"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/9seconds/httransform/v2"
	"github.com/9seconds/httransform/v2/layers"
	_ "github.com/mattn/go-sqlite3"
	l "github.com/rubiojr/eyez/internal/layers"
)

func main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range signals {
			cancel()
		}
	}()

	caCert, err := ioutil.ReadFile("certs/rootCA.crt")
	if err != nil {
		panic(err)
	}
	caPrivateKey, _ := ioutil.ReadFile("certs/rootCA.key")

	persistance, err := l.NewPersist()
	if err != nil {
		panic(err)
	}

	opts := httransform.ServerOpts{
		TLSCertCA:     caCert,
		TLSPrivateKey: caPrivateKey,
		Layers: []layers.Layer{
			persistance,
			l.Stdout{},
			layers.TimeoutLayer{
				Timeout: 3 * time.Minute,
			},
		},
	}

	proxy, err := httransform.NewServer(ctx, opts)
	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", ":1080")
	if err != nil {
		panic(err)
	}

	if err := proxy.Serve(listener); err != nil {
		panic(err)
	}
}
