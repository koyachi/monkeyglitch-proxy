package main

import (
	"github.com/codegangsta/cli"
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/image"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type MonkeyGlitcher struct {
	id  string
	R   io.ReadCloser
	nth int
}

func (m *MonkeyGlitcher) Read(p []byte) (n int, err error) {
	n, err = m.R.Read(p)
	i := rand.Float64()
	if i > 0.2 {
		index := int(float64(n) * i)
		v := rand.Float64()
		p[index] = 0xFF * byte(v)
	}
	m.nth = m.nth + 1
	return
}

func (m MonkeyGlitcher) Close() error {
	return m.R.Close()
}

func main() {
	app := cli.NewApp()
	app.Name = "monkeyglitch-proxy"
	app.Usage = "A HTTP proxy for corrupting images."
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.IntFlag{"port, p", 8080, "TCP/IP port number"},
	}
	app.Action = func(c *cli.Context) {
		proxy := goproxy.NewProxyHttpServer()
		proxy.OnResponse(goproxy_image.RespIsImage).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			resp.Body = &MonkeyGlitcher{ctx.Req.URL.String(), resp.Body, 0}
			return resp
		})
		port := c.Int("port")
		portStr := ":" + strconv.Itoa(port)
		log.Printf("Starting %vserver with port %v\n", app.Name, portStr)
		log.Fatal(http.ListenAndServe(portStr, proxy))
	}
	app.Run(os.Args)

}
