package examplev1

import (
	"context"
	"fmt"

	"github.com/anzboi/proto-playground/pkg/api/example"
)

type Impl struct{}

func (Impl) SayHello(ctx context.Context, req *example.HelloRequest) (*example.HelloResponse, error) {
	name := "world"
	if req.GetName() != "" {
		name = req.GetName()
	}
	response := fmt.Sprintf("Hello %s", name)
	return &example.HelloResponse{Message: response}, nil
}
