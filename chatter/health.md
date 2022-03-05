# GRPC health

Here we discuss the standard grpc healthcheck service and how to test against it with chatter.

## Overview

GRPC defines a health check service [here](https://github.com/grpc/grpc/blob/master/src/proto/grpc/health/v1/health.proto). This service is rapidly becoming standard, as of december 2021, [k8s added their grpc health probe as a built in feature](https://kubernetes.io/blog/2018/10/01/health-checking-grpc-servers-on-kubernetes/).

The health service defines 2 rpcs, a simple `Check` RPC and a `Watch` RPC. Both of these OPTIONALLY take a grpc `service` name as input, and returns whether or not the server is currently serving that service. If no service is given, the result should be interpreted the health of the server as a whole.

## Testing

The chatter service implements a simple health check server. It does not actively monitor its own health, so status will not change over the services lifespan.

To change the health status for testing purposes, set the `HEALTH` environment variables to one of the following enum values (`SERVING`, `NOT_SERVING`, `SERVICE_UNKNOWN`)