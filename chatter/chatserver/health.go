package main

import (
	"context"

	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Basic, static health server implementation
type healthServer struct {
	health.UnimplementedHealthServer
	services map[string]health.HealthCheckResponse_ServingStatus
}

// Check implements grpc.health.v1.Health.Check
func (h *healthServer) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	if req.GetService() == "" {

	}
	currentHealth, ok := h.services[req.GetService()]
	if !ok {
		return nil, status.Error(codes.NotFound, "unknown service")
	}
	return &health.HealthCheckResponse{Status: currentHealth}, nil
}

// Watch implements grpc.health.v1.Health.Watch
//
// This basic implementation does not watch for status changes, so the first
// response will be the only response sent, but stream will remain open to satisfy
// the RPC contract.
func (h *healthServer) Watch(req *health.HealthCheckRequest, svr health.Health_WatchServer) error {
	currentHealth, ok := h.services[req.GetService()]
	if !ok {
		svr.Send(&health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVICE_UNKNOWN})
	}
	svr.Send(&health.HealthCheckResponse{Status: currentHealth})

	// Wait until the client closes the stream
	<-svr.Context().Done()
	return nil
}
