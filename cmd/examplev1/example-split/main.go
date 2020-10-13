package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/anzboi/proto-playground/cmd/examplev1"
	"github.com/anzboi/proto-playground/pkg/api/example"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	stdoutexporter "go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	exp = flag.String("otel_exporter", "stdout", "select an opentelemetry exporter to send traces to")
)

func setupOpenTelemetry(ctx context.Context) (func(), error) {
	close := func() {}
	var exporter interface {
		trace.SpanExporter
		metric.Exporter
	}
	var err error

	resource, err := resource.Detect(context.Background(), &gcp.GKE{})
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(*exp) {
	case "", "stdout":
		exporter, err = stdoutexporter.NewExporter(stdoutexporter.WithPrettyPrint())
	case "coll", "colector":
		exporter, err = otlp.NewExporter()
	}
	if err != nil {
		return nil, err
	}
	close = func() { defer exporter.Shutdown(ctx); close() }

	pusher := push.New(basic.New(simple.NewWithInexpensiveDistribution(), metric.CumulativeExporter), exporter,
		push.WithResource(resource),
	)
	pusher.Start()
	close = func() { defer pusher.Stop(); close() }

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithResource(resource),
		sdktrace.WithSyncer(exporter),
	)

	global.SetTracerProvider(tp)
	global.SetMeterProvider(pusher.MeterProvider())
	return close, nil
}

func main() {
	flag.Parse()

	closeOtel, err := setupOpenTelemetry(context.Background())
	if err != nil {
		panic(err)
	}
	defer closeOtel()

	errCh := make(chan error)
	go func() {
		errCh <- RunGRPCServer(examplev1.Impl{})
	}()
	go func() {
		errCh <- RunHTTPGateway()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal("runtime error, shutting down:", err)
	}
}

func RunGRPCServer(impl example.HelloWorldServer) error {
	svr := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(global.Tracer("")),
			// other interceptors
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(global.Tracer("")),
			// other interceptors
		),
	)
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
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(global.Tracer("")),
			// other interceptors
		),
		grpc.WithChainStreamInterceptor(
			otelgrpc.StreamClientInterceptor(global.Tracer("")),
			// other interceptors
		),
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
