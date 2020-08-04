package main

import (
	"context"
	"fmt"
	"net"

	"github.com/anzboi/proto-playground/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type impl struct {
	rpc.UnimplementedCatalogServer
	rpc.UnimplementedChatRoomServer

	rooms map[string]*Topic
}

type Topic struct {
	messages    []*rpc.Message
	subscribers map[string]chan *rpc.Message
}

func (t *Topic) Push(pusher string, m *rpc.Message) {
	t.messages = append(t.messages, m)
	for name, sub := range t.subscribers {
		if name != pusher {
			sub <- m
		}
	}
}

func (t *Topic) Subscribe(name string, ch chan *rpc.Message) {
	t.subscribers[name] = ch
}

func (t *Topic) Unsubscribe(name string) {
	delete(t.subscribers, name)
}

func (i *impl) ListProducts(context.Context, *rpc.ListProductsRequest) (*rpc.ProductList, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (i *impl) GetProduct(context.Context, *rpc.GetProductRequest) (*rpc.Product, error) {
	return &rpc.Product{
		ProductId: 123,
		Name:      "Rice Cooker",
	}, nil
}

func (i *impl) Chat(stream rpc.ChatRoom_ChatServer) error {
	first, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, "Failed to read initial message")
	}

	name := first.GetJoin().GetName()
	var room *Topic
	ch := make(chan *rpc.Message)
	if topic, ok := i.rooms[first.GetJoin().GetRoomId()]; ok {
		room = topic
	} else {
		room = &Topic{
			subscribers: map[string]chan *rpc.Message{},
		}
	}
	room.Subscribe(name, ch)

	for _, msg := range room.messages {
		stream.Send(msg)
	}

	go func() {
		select {
		case m := <-ch:
			stream.Send(m)
		case <-stream.Context().Done():
			room.Unsubscribe(name)
		}
	}()

	go func() {
		for {
			m, _ := stream.Recv()
			room.Push(name, m.GetMessage())
		}
	}()

	<-stream.Context().Done()
	return nil
}

func main() {
	svr := grpc.NewServer()
	rpc.RegisterCatalogServer(svr, &impl{})
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
