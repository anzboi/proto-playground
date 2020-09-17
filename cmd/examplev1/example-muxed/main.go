package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/anzboi/proto-playground/cmd/examplev1"
	"github.com/anzboi/proto-playground/pkg/api/example"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// setup server
	svr := grpc.NewServer()
	example.RegisterHelloWorldServer(svr, examplev1.Impl{})
	reflection.Register(svr)

	// setup net listener
	addr := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		addr = ":" + envPort
	}
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// setup http gateway
	mux := runtime.NewServeMux()
	example.RegisterHelloWorldHandlerFromEndpoint(context.Background(), mux, "localhost"+addr, []grpc.DialOption{
		grpc.WithInsecure(),
	})

	entrypoint := grpcDispatcher(context.Background(), svr, mux)

	// run
	log.Printf("Listening on %s", addr)
	if err := http.Serve(lis, entrypoint); err != nil {
		panic(err)
	}
}

// Use x/net/http2/h2c so we can have http2 cleartext connections. The default
// Go http server does not support it. We also cannot plug into the grpc
// http2 server.
// From: https://github.com/philips/grpc-gateway-example/issues/22#issuecomment-490733965
func grpcDispatcher(ctx context.Context, grpcHandler http.Handler, httpHandler http.Handler) http.Handler {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			log.Print("dispatch: grpc")
			grpcHandler.ServeHTTP(w, r)
		} else {
			log.Print("dispatch: http gateway")
			httpHandler.ServeHTTP(w, r)
		}
	})
	return h2c.NewHandler(hf, &http2.Server{})
}
