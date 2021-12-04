package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/9seconds/httransform/v2"
	"github.com/9seconds/httransform/v2/layers"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type StdoutLayer struct{}

func main() {
	var err error
	db, err = initDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// For demo purpose we are going to close by SIGINT and SIGTERM
	// signals.
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
	opts := httransform.ServerOpts{
		TLSCertCA:     caCert,
		TLSPrivateKey: caPrivateKey,
		Layers: []layers.Layer{
			StdoutLayer{},
			layers.TimeoutLayer{
				Timeout: 3 * time.Minute,
			},
		},
	}

	proxy, err := httransform.NewServer(ctx, opts)
	if err != nil {
		panic(err)
	}

	// We bind our proxy to the port 3128 and all interfaces.
	listener, err := net.Listen("tcp", ":1080")
	if err != nil {
		panic(err)
	}

	if err := proxy.Serve(listener); err != nil {
		panic(err)
	}
}

func (StdoutLayer) OnRequest(ctx *layers.Context) error {
	time.Sleep(1 * time.Second)
	req := ctx.Request()
	fmt.Println()
	fmt.Println(urlStyle.Render(req.URI().String()))
	fmt.Printf("%s: %d\n", keyStyle.Render("Status"), ctx.Response().StatusCode())
	fmt.Printf("%s: %s\n", keyStyle.Render("Method"), string(req.Header.Method()))
	fmt.Printf("%s: %s\n", keyStyle.Render("Path"), string(req.URI().Path()))
	fmt.Println(keyStyle.Render("Headers:"))
	for _, line := range strings.Split(strings.TrimSuffix(string(req.Header.RawHeaders()), "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "Authorization:") {
			fmt.Println(headersStyle.Render("Autorization: *****"))
		} else {
			fmt.Println(headersStyle.Render(line))
		}
	}
	body := req.Body()
	if len(body) > 0 {
		fmt.Printf("%s: %d bytes\n", keyStyle.Render("Body Size"), len(body))
		if len(body) > 8192 {
			fmt.Printf("%s: %s\n", keyStyle.Render("Body"), "[too large]")
		} else {
			fmt.Printf("%s: %s\n", keyStyle.Render("Body"), (req.Body()))
		}
	} else {
		fmt.Printf("%s: %s\n", keyStyle.Render("Body"), "N/A")
	}
	_, err := db.Exec("INSERT INTO "+defaultCaptureCollection+" (uuid, url, body, path, headers, date_end, status, method) VALUES (?,?,?,?,?,?,?,?)",
		uuid.New().String(),
		req.URI().String(),
		body,
		req.URI().Path(),
		req.Header.RawHeaders(),
		time.Now(),
		ctx.Response().StatusCode(),
		req.Header.Method(),
	)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (StdoutLayer) OnResponse(ctx *layers.Context, err error) error {
	return nil
}
