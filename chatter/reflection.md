# Reflection

Here we give a brief overview of descriptors, reflection, and their place in the world of protos and grpc.

## Descriptors

TODO

## The reflection API

TODO

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
