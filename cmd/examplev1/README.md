# ExampleV1

ExampleV1 is a set of basic microservices that all implement the same gRPC service with a few different configurations

## Run

All the examples can be run using `go run` or by building a docker container with the relevant dockerfile. With a service running, you can hit is with the following requests

```sh
# grpcurl the service
grpcurl --plaintext -d '{"name": "anzboi"}' <host>:<port> anzboi.example.v1/HelloWorld

# curl the gateway
curl http://<host>:<port>/v1/example/hello?name=anzboi
```

### example-cmuxed

```sh
# go run
go run example-cmuxed

# docker
docker build -t example-cmuxed -f cmuxed.dockerfile ../../
docker run -it -p 8080:8080 example-cmuxed
```

* cmux for muxing http vs grpc requests. Serves both on port 8080 (settable via env var `PORT`)
* external calls for gateway -> grpc connection
* No tracing/metrics instrumentation

### example-muxed

```sh
# go run
go run example-muxed

# docker
docker build -t example-muxed -f muxed.dockerfile ../../
docker run -it -p 8080:8080 example-muxed
```

* http mux for muxing http vs grpc requests. Serves both on port 8080 (settable via env var `PORT`)
* external calls for gateway -> grpc connection
* No tracing/metrics instrumentation

### example-split

```sh
# go run
go run example-split

# docker
docker build -t example-split -f split.dockerfile ../../
docker run -it -p 8080:8080 example-split
```

* http mux for muxing http vs grpc requests. Serves grpc on port 8080 (`PORT`) and http on 8081 (`HTTP_PORT`)
* external calls for gateway -> grpc connection
* Opentelemetry trace and metrics instrumentation

## Open Telemetry

Open telemetry is instrumented in the example-split microservice. Use the `otel_exporter` flag to set the exporter you wish to test with

```sh
# stdout
example-split --otel_exporter=stdout

# opentelemetry collector
example-split --otel_exporter=collector
```

## Exposing an HTTP-REST gateway

There are a few ways to do this related to whether you want to serve both grpc and http on the same port or not, and how you connect the http gateway to grpc

### Serving on a Port

**cmuxed**. [cmux](https://github.com/soheilhy/cmux) is a connection multiplexer that is able to mux a conn object. This enables exposing the servers on two different listeners that are connected to the same real conn object via cmux.

Advantage: multiplexing the conn means just about any tool tat works with a conn will work with `cmux`. In practice, this means the gRPC server can use the `Serve(net.Listener)` method. Can serve all endpoints on one port.

Disdvantage: Tends to not play well with reverse proxies sitting in front of your service.

**muxed** (or http-muxed). The http mux pattern involves using an `http2.Server` and dispatching requests to a handler based on request information like `Proto` and `Content-Type` header. `Proto=2` and `Content-Type=application/grpc` is enough to distinguish between an http and a grpc request, and the http2 server is compatible with both proto versions.

Advantage: Tends to play better with reverse proxies. Can server all endpoints on one port

Disadvantage: You backend server MUST expose a `ServeHTTP` method. This is unusual for gRPC and third part tools that expect the `grpcServer.Server(net.Listener)` method will not work

**split**. This involves service both http and grpc on different ports.

Advantage: Complete flexibility in instrumentation, requests do not need to be muxed in any form

Disadvantage: Having to manage two ports. This complicates deployments and networking rules when considering the larger picture.

### Connecting Gateways to gRPC services

GRPC gateways perform a very simple task. Convert an http request to a gRPC one, hit the gRPC endpoint, and convert the response back to http. There are essentially two ways to hit the gRPC endpoint

**Direct**. This refers to sending an http request directly to the implementer of a grpc service. This completely bypasses anything

Advantages: Trace and metric information are not doubled up. More efficient (request does not have to go back out to the local network).

Disadvantage: Bypasses grpc interceptor stack, forcing you to create and maintain two interceptor stacks. Brings in the possibility for differences between the two

**External**. This refers to connecting to the grpc service by send a request back out as if it were connecting to any other downstream. The request comes back into the service as if it were an external request.

Advantages: all requests are forced to go through the gRPC interceptor stack (consistency, easier to maintain).

Disadvantages: innefficient (requires an extra network hop). tracing and metrics information is sort of doubled up (you could argue this is correct but it seems superfluous).
