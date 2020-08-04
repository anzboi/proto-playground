# Chatroom

A gRPC chat room service. Create, Delete chat rooms and connect to them to chat with others.

This service serves to experiment with bidirectional streams. You can start the chat service and inspect http2 (grpc) packets on your system with a packet viewer like [wireshark](https://www.wireshark.org/).

- [protos](../proto/rpc/chat.proto)
- [server](../cmd/chatserver)
- [chatter client](../cmd/chatter)

## Run

```sh
# Run the server
$ go run ./cmd/chatserver

# Run a client
$ go run ./cmd/chatter
```

## Creating, Listing and Deleting rooms

A chat service admin can create chat rooms Using the `CreateChatRoom` RPC

```sh
$ grpcurl host:port rpc.ChatService/CreateChatRoom
{
    "room_id": "126859",
}
```

Existing rooms can be listed using `ListChatRooms`, Deleting rooms can be done using `DeleteChatRoom`

## Using the chatter client

The `chatter` is the client CLI to conduct chat. Think of this as a basic client that can consume the chat service

### Joining a chat room

The tool first connects to the server and uses the ListRooms endpoint to discover active chat rooms

```sh
$ go run ./cmd/chatter
Chat rooms:
[1]: 1234
[2]: 5678
...
```

The tool will then prompt you for a room the join, and a name tag so other can identify you

```sh
Pick a room: 2
Enter a Name Tag: Johnny Boi
```

The tool will then join the chat by sending an initial request to the chat server asking to join a particular room with the given name tag.

### Handling chat

After the initial join request, the tool will start two processes.

1. Any message recieved by the server will be displayed to stdout with the users name tag
1. Any message entered on stdin by the user will be sent to the server

This will typically look something like

```sh
Joining room
Bob: Has entered the chat
Bob: Hey Johnny
Hey bob                     # user entered message
Alice: Has entered the chat
Alice: Greetings!
```

### Deleted rooms

Rooms can be deleted while users are still connected to them. The server will issue a final message as coming from the chat server stating that the room is close, then terminate the connection.

```sh
chat server: chat room is closing                # Final message
message recieved but could not be displayed: EOF # client error, this could probably be improved to a proper status
```
