package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type grpcServer struct {
	UnimplementedEchoServer
}

func (grpcServer) Echo(
	ctx context.Context,
	req *EchoRequest,
) (*EchoResponse, error) {
	return &EchoResponse{
		Message: req.Message,
	}, nil
}

type httpServer struct{}

func (httpServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hello world!")
}

func main() {
	g := grpc.NewServer()
	RegisterEchoServer(g, grpcServer{})

	h := httpServer{}

	serve := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.ProtoMajor == 2 && strings.HasPrefix(
			req.Header.Get("Content-Type"), "application/grpc",
		) {
			g.ServeHTTP(res, req)
		} else {
			h.ServeHTTP(res, req)
		}
	})

	ch := make(chan error)

	go func() {
		// Broken
		ch <- http.ListenAndServe(":8080", serve)
	}()
	go func() {
		// Fixed
		ch <- http.ListenAndServe(":8081",
			h2c.NewHandler(serve, &http2.Server{}))
	}()
	log.Println(<-ch)
}
