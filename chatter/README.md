# Chatter

Chatter is a very simple chat server over gRPC. Chatter is used to demonstrate and test various features and tooling against gRPC. It exposes a few RPCs including unary RPCs, as well as client, server and bidirectional streams.

## Install

You will need to [install go](https://golang.org/doc/install) in order to run chatter

You can install the chat client and server with

```sh
go install ./...
```

Alternatively the server can be built and run in a docker container using

```sh
docker build -t chatserver -f Dockerfile ..
docker run -d -p 8080:8080 chatserver
```

## Run

### Run the server

```sh
go run ./chatserver
```

### Create a chat room

```sh
# add --insecure for no tls
chatter --create
```

### Run the client

```sh
# add --insecure for no tls
go run ./chatter
```

The chatter client will retrieve a list of rooms and ask which one you would like to join, followed by a user name

```sh
Chat rooms:
[1]: 893f1e077e88c01341a8
Pick a room: 1
Enter Name tag: Me
Me: Has joined the chat

```

Run multiple clients to simulate multiple chatting users. Each connected client opens their own bidirectional stream to the server.

## Behind the scenes

The chat server implements a really crufty pub-sub model and allows goroutines to subscribe to a chatroom and publish messages to it. The Chat RPC simply connects a grpc stream to a chat room.

When the server receives a stream message from the client, it publishes it to the channel so it is picked up by all the subscribers. When the stream picks up a message on the chat room, it streams it back to the client.

On top of this, the chat server to host multiple chat rooms which can be created and deleted through the `CreateChatRoom` and `DeleteChatRoom` RPCs respectively. The first message the server receives from the client must contain JoinParameters which tell it which chat room the stream wants to connect to, and the user tag the stream will be posting messages as. All subsequent messages will be treated as ordinary chat messages.

## Invoking RPCs from descriptors

The `chatter-from-descriptors` executable demonstrates the typical way you would invoke RPCs using proto descriptors instead of relying on generated code.

The basic sequence of operations is as follows

1. Obtain a set of descriptors from some kind of descriptor provider. Providers could include...
    1. GRPC reflection (runs directly on the server, if configured to do so)
    2. proto files (sourced locally or from a file server)
    3. descriptor sets (binary encoding of proto files that serve the same job as the files themselves).
2. Search the descriptors for a **method descriptor** matching the endpoint you wish to invoke.
3. Create request/response objects from the descriptors and some input data.
4. Invoke the method

You can run this sequence of operations against a running chat server using...

```sh
./chatter-from-descriptors --rpc rpc.ChatServer.CreateChatRoom --data='{"room_id":"abcdefg"}'
```

NOTE: this executable is fixed to fetch descriptors from grpc reflection running on the server, and does not support streaming.
