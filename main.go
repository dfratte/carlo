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

	ukEast := createBackendServer(ctx, "3333", mux)
	ukWest := createBackendServer(ctx, "3334", mux)
	france := createBackendServer(ctx, "4444", mux)
	germany := createBackendServer(ctx, "5555", mux)

	go func() {
		err := ukEast.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("ukEast closed\n")
		} else if err != nil {
			fmt.Printf("error listening for ukEast: %s\n", err)
		}
		cancelCtx()
	}()
	go func() {
		err := ukWest.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("ukWest closed\n")
		} else if err != nil {
			fmt.Printf("error listening for ukWest: %s\n", err)
		}
		cancelCtx()
	}()
	go func() {
		err := france.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("france closed\n")
		} else if err != nil {
			fmt.Printf("error listening for france: %s\n", err)
		}
		cancelCtx()
	}()
	go func() {
		err := germany.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("germany closed\n")
		} else if err != nil {
			fmt.Printf("error listening for germany three: %s\n", err)
		}
		cancelCtx()
	}()

	ukServers := []server{
		CreateBackend("http://localhost:3333"),
		CreateBackend("http://localhost:3334"),
	}

	franceServers := []server{
		CreateBackend("http://localhost:4444"),
	}

	germanyServers := []server{
		CreateBackend("http://localhost:5555"),
	}

	regions := []*region{
		CreateRegion(ukServers, 0, "uksouth"),
		CreateRegion(franceServers, 0, "francecentral"),
		CreateRegion(germanyServers, 0, "germanywestcentral"),
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
