# Proto playground

A proto and grpc playground. This repo contains a bunch of example proto and grpc definitions and implementations to play around with and experiment.

- [Chatter](./chatter). A Chat room service with ability to create and manage separate chat rooms, and connect and chat with other users. Exposes unary, vlient, server and bidirectional stream.
  - includes demonstration of using descriptors to discover proto specs dynamically and call endpoints.
- [Example](./cmd/examples). Demostrates a few ways to serve gRPC and http APIs in a single container (multi-port, a few single-port implementations).