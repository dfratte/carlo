package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

const maxLatency int = 100

type server interface {
	Address() string
	Latency() int
	Serve(rw http.ResponseWriter, r *http.Request)
}

type backend struct {
	address string
	latency int
	proxy   *httputil.ReverseProxy
}

func CreateBackend(a string) *backend {
	backendURL, _ := url.Parse(a)
	return &backend{
		address: a,
		latency: generateRandomInt(maxLatency),
		proxy:   httputil.NewSingleHostReverseProxy(backendURL),
	}
}

func (b *backend) Address() string {
	return b.address
}

func (b *backend) Latency() int {
	return b.latency
}

func (b *backend) Serve(rw http.ResponseWriter, r *http.Request) {
	b.latency = generateRandomInt(maxLatency)
	b.proxy.ServeHTTP(rw, r)
}
