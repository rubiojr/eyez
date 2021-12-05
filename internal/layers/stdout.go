package layers

import (
	"fmt"
	"time"

	"github.com/9seconds/httransform/v2/layers"
	"github.com/rubiojr/eyez/internal/styles"
)

type Stdout struct{}

func (Stdout) OnRequest(ctx *layers.Context) error {
	time.Sleep(1 * time.Second)
	req := ctx.Request()
	fmt.Println()
	fmt.Println(styles.Url.Render(req.URI().String()))
	fmt.Printf("%s: %d\n", styles.Key.Render("Status"), ctx.Response().StatusCode())
	fmt.Printf("%s: %s\n", styles.Key.Render("Method"), string(req.Header.Method()))
	fmt.Printf("%s: %s\n", styles.Key.Render("Path"), string(req.URI().Path()))
	fmt.Println(styles.Key.Render("Headers:"))
	req.Header.VisitAll(func(key, value []byte) {
		h := fmt.Sprintf("%s: %s", key, value)
		if string(key) == "Authorization" {
			h = fmt.Sprintf("%s: %s", key, "[REDACTED]")
		}
		fmt.Println(styles.Header.Render(h))
	})
	body := req.Body()
	if len(body) > 0 {
		fmt.Printf("%s: %d bytes\n", styles.Key.Render("Body Size"), len(body))
		if len(body) > 8192 {
			fmt.Printf("%s: %s\n", styles.Key.Render("Body"), "[too large]")
		} else {
			fmt.Printf("%s: %s\n", styles.Key.Render("Body"), (req.Body()))
		}
	} else {
		fmt.Printf("%s: %s\n", styles.Key.Render("Body"), "N/A")
	}
	fmt.Printf(styles.Tag.Render("connect") + " " + styles.Tag.Render("core"))
	fmt.Println()

	return nil
}

func (Stdout) OnResponse(ctx *layers.Context, err error) error {
	return nil
}
