# Chatter

Chatter is a very simple chat server over gRPC. Chatter is intended to be used to test bidirectional streams.

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
