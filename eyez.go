package eyez

import (
	"context"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	_ "embed"

	"github.com/9seconds/httransform/v2"
	"github.com/9seconds/httransform/v2/layers"
	l "github.com/rubiojr/eyez/internal/layers"
)

//go:embed certs/rootCA.crt
var caCert []byte

//go:embed certs/rootCA.key
var caKey []byte

type Proxy struct {
	p    *httransform.Server
	opts *ProxyOptions
}

type ProxyOptions struct {
	Port   int
	BindIP string
	Layers []layers.Layer
	CACert string
	CAKey  string
}

func New(ctx context.Context, opts *ProxyOptions) (*Proxy, error) {
	var err error

	caCertPath := opts.CACert
	if caCertPath != "" {
		caCert, err = ioutil.ReadFile(caCertPath)
		if err != nil {
			return nil, err
		}
	}

	caKeyPath := opts.CAKey
	if caKeyPath != "" {
		caKey, err = ioutil.ReadFile(caKeyPath)
		if err != nil {
			return nil, err
		}
	}

	if opts.Port == 0 {
		opts.Port = 1080
	}

	if opts.Layers == nil {
		if opts.Layers, err = DefaultLayers(); err != nil {
			return nil, err
		}
	}

	popts := httransform.ServerOpts{
		TLSCertCA:     caCert,
		TLSPrivateKey: caKey,
		Layers:        opts.Layers,
	}

	proxy := Proxy{opts: opts}
	proxy.p, err = httransform.NewServer(ctx, popts)

	return &proxy, err
}

func DefaultLayers() ([]layers.Layer, error) {
	persistance, err := l.NewPersist()
	if err != nil {
		return nil, err
	}
	stdout := l.Stdout{}
	timeout := layers.TimeoutLayer{
		Timeout: 3 * time.Minute,
	}

	return []layers.Layer{persistance, stdout, timeout}, nil
}

func (p *Proxy) Serve() error {
	bip := p.opts.BindIP
	port := strconv.Itoa(p.opts.Port)

	listener, err := net.Listen("tcp", bip+":"+port)
	if err != nil {
		return err
	}

	return p.p.Serve(listener)
}
