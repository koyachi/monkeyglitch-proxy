package main

import (
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/image"
	"io"
	"log"
	"math/rand"
	"net/http"
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
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnResponse(goproxy_image.RespIsImage).DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp.Body = &MonkeyGlitcher{ctx.Req.URL.String(), resp.Body, 0}
		return resp
	})
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
