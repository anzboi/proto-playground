package main

import (
	"context"

	"github.com/anzboi/proto-playground/pkg/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogImpl struct {
	rpc.UnimplementedCatalogServer
}

// Implementations for two of the catalog RPCs
func (i *CatalogImpl) ListProducts(context.Context, *rpc.ListProductsRequest) (*rpc.ProductList, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (i *CatalogImpl) GetProduct(context.Context, *rpc.GetProductRequest) (*rpc.Product, error) {
	return &rpc.Product{
		ProductId: 123,
		Name:      "Rice Cooker",
	}, nil
}
