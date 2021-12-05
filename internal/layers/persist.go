package layers

import (
	"fmt"
	"time"

	"github.com/9seconds/httransform/v2/layers"
	"github.com/google/uuid"
	"github.com/rubiojr/eyez/internal/db"
)

type Persist struct{}

func NewPersist() (layers.Layer, error) {
	err := db.InitDB()
	if err != nil {
		return nil, err
	}

	return &Persist{}, nil
}

func (Persist) OnRequest(ctx *layers.Context) error {
	req := ctx.Request()
	body := req.Body()
	_, err := db.Exec("INSERT INTO "+db.DefaultCaptureCollection+" (uuid, url, body, path, headers, date_end, status, method) VALUES (?,?,?,?,?,?,?,?)",
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

func (Persist) OnResponse(ctx *layers.Context, err error) error {
	return nil
}
