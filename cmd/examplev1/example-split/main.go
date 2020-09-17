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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	errCh := make(chan error)
	go func() {
		errCh <- RunGRPCServer(examplev1.Impl{})
	}()
	go func() {
		errCh <- RunHTTPGateway()
	}()

	err := <-errCh
	log.Fatal("runtime error, shutting down:", err)
}

func RunGRPCServer(impl example.HelloWorldServer) error {
	svr := grpc.NewServer()
	example.RegisterHelloWorldServer(svr, impl)
	reflection.Register(svr)

	lis, err := net.Listen("tcp", grpcAddr())
	if err != nil {
		return err
	}

	log.Printf("grpc listening on %s", grpcAddr())
	return svr.Serve(lis)
}

func RunHTTPGateway() error {
	mux := runtime.NewServeMux()
	example.RegisterHelloWorldHandlerFromEndpoint(context.Background(), mux, "localhost"+grpcAddr(), []grpc.DialOption{
		grpc.WithInsecure(),
	})

	log.Printf("http listening on %s", httpAddr())
	return http.ListenAndServe(httpAddr(), mux)
}

func grpcAddr() string {
	addr := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		addr = ":" + envPort
	}
	return addr
}

func httpAddr() string {
	addr := ":8081"
	if envPort := os.Getenv("HTTP_PORT"); envPort != "" {
		addr = ":" + envPort
	}
	return addr
}
