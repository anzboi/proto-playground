package main

import (
	"fmt"
	"net"

	"github.com/anzboi/proto-playground/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	svr := grpc.NewServer()
	rpc.RegisterCatalogServer(svr, &CatalogImpl{})
	rpc.RegisterChatServiceServer(svr, NewChatService())
	reflection.Register(svr)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on 8080")
	if err := svr.Serve(lis); err != nil {
		panic(err)
	}

	grpc.Dial("localhost:8080")

	grpc.Dial("localhost:8080", grpc.WithInsecure())

}
