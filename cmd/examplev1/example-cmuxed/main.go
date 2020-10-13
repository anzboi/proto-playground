package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/anzboi/proto-playground/cmd/examplev1"
	"github.com/anzboi/proto-playground/pkg/api/example"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// setup net listener
	addr := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		addr = ":" + envPort
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// setup cmux
	m := cmux.New(lis)
	grpcLis := m.Match(cmux.HTTP2(), cmux.HTTP2HeaderFieldPrefix("content-type", "appliction/grpc"))
	httpLis := m.Match(cmux.Any())

	// setup server
	svr := grpc.NewServer()
	example.RegisterHelloWorldServer(svr, examplev1.Impl{})
	reflection.Register(svr)

	// setup http gateway
	gateway := runtime.NewServeMux()
	example.RegisterHelloWorldHandlerFromEndpoint(context.Background(), gateway, "localhost"+addr, []grpc.DialOption{
		grpc.WithInsecure(),
	})

	errCh := make(chan error)

	// Start servers
	go func() {
		log.Printf("listening on %s", addr)
		errCh <- m.Serve()
	}()
	go func() {
		errCh <- svr.Serve(grpcLis)
	}()
	go func() {
		errCh <- http.Serve(httpLis, gateway)
	}()

	if err := <-errCh; err != nil {
		panic(err)
	}
}
