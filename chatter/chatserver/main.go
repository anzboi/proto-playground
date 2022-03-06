package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/anzboi/proto-playground/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	svr := grpc.NewServer()
	rpc.RegisterCatalogServer(svr, &CatalogImpl{})
	rpc.RegisterChatServiceServer(svr, NewChatService())

	// Health server
	// you may test different health responses by setting values here
	health := getHealth()
	healthServer := &healthServer{services: map[string]grpc_health_v1.HealthCheckResponse_ServingStatus{
		// empty string is the default catch-all
		"":                health,
		"rpc.Catalog":     health,
		"rpc.ChatService": health,
		"grpc.reflection.v1alpha.ServerReflection": health,
		"grpc.health.v1.Health":                    health,
	}}
	grpc_health_v1.RegisterHealthServer(svr, healthServer)
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

func getHealth() grpc_health_v1.HealthCheckResponse_ServingStatus {
	health := os.Getenv("HEALTH")
	health = strings.ToLower(health)
	switch health {
	case "", "1", "serving":
		return grpc_health_v1.HealthCheckResponse_SERVING
	case "2", "not_serving":
		return grpc_health_v1.HealthCheckResponse_NOT_SERVING
	case "3", "service_unknown":
		return grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN
	default:
		return grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
}
