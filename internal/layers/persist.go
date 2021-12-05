package layers

import (
	"fmt"
	"strings"
	"time"

	"github.com/9seconds/httransform/v2/layers"
	"github.com/google/uuid"
	"github.com/rubiojr/eyez/internal/db"
)

type Persist struct{}

func NewPersist(path string) (layers.Layer, error) {
	err := db.InitDB(path)
	if err != nil {
		return nil, err
	}

	return &Persist{}, nil
}

func (Persist) OnRequest(ctx *layers.Context) error {
	req := ctx.Request()
	var sb strings.Builder
	req.Header.VisitAll(func(key, value []byte) {
		if string(key) == "Authorization" {
			sb.WriteString("Authorization: [REDACTED]\n")
			return
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	})
	body := req.Body()
	_, err := db.Exec("INSERT INTO "+db.DefaultCaptureCollection+" (uuid, url, body, path, headers, date_end, status, method) VALUES (?,?,?,?,?,?,?,?)",
		uuid.New().String(),
		req.URI().String(),
		body,
		req.URI().Path(),
		sb.String(),
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
