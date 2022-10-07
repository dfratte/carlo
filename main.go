package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := ctx.Value(keyServerAddr)
	fmt.Fprintf(w, "Hello from %s! This is a carbon-aware server response\n", s)
}

func createBackendServer(ctx context.Context, port string, mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)

	ctx, cancelCtx := context.WithCancel(context.Background())

	irelandEast := createBackendServer(ctx, "3333", mux)
	irelandWest := createBackendServer(ctx, "3334", mux)
	london := createBackendServer(ctx, "4444", mux)
	paris := createBackendServer(ctx, "5555", mux)

	go func() {
		err := irelandEast.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("irelandEast closed\n")
		} else if err != nil {
			fmt.Printf("error listening for irelandEast: %s\n", err)
		}
		cancelCtx()
	}()
	go func() {
		err := irelandWest.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("irelandWest closed\n")
		} else if err != nil {
			fmt.Printf("error listening for irelandWest: %s\n", err)
		}
		cancelCtx()
	}()
	go func() {
		err := london.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("london closed\n")
		} else if err != nil {
			fmt.Printf("error listening for london: %s\n", err)
		}
		cancelCtx()
	}()
	go func() {
		err := paris.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("paris closed\n")
		} else if err != nil {
			fmt.Printf("error listening for paris three: %s\n", err)
		}
		cancelCtx()
	}()

	irelandServers := []server{
		CreateBackend("http://localhost:3333"),
		CreateBackend("http://localhost:3334"),
	}

	parisServers := []server{
		CreateBackend("http://localhost:4444"),
	}

	londonServers := []server{
		CreateBackend("http://localhost:5555"),
	}

	regions := []*region{
		CreateRegion(irelandServers, 0, "eu-west-1"),
		CreateRegion(parisServers, 0, "eu-west-2"),
		CreateRegion(londonServers, 0, "eu-west-3"),
	}

	lb := NewLoadBalancer(regions)
	handler := func(rw http.ResponseWriter, r *http.Request) {
		lb.HandleRequest(rw, r)
	}

	http.HandleFunc("/", handler)
	fmt.Printf("carlo load balancer listening at 'localhost:%s'\n", "8000")
	http.ListenAndServe(":"+"8000", nil)

	<-ctx.Done()
}
