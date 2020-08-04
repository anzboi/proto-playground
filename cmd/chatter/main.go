package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/anzboi/proto-playground/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	roomID   = flag.String("room", "", "room id of chat room to connect to")
	name     = flag.String("name", "anonymous", "name tag to use in chat room")
	host     = flag.String("host", "localhost:8080", "chat server host address")
	insecure = flag.Bool("insecure", false, "set insecure to true to use http instead of https")
)

func main() {
	flag.Parse()
	opts := []grpc.DialOption{}
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		tc := credentials.NewTLS(nil)
		opts = append(opts, grpc.WithTransportCredentials(tc))
	}
	cc, err := grpc.Dial(*host, opts...)
	if err != nil {
		panic(err)
	}

	client, err := rpc.NewChatServiceClient(cc).Chat(context.Background())
	if err != nil {
		panic(err)
	}

	join := &rpc.ChatMessage{
		Join: &rpc.JoinChat{
			RoomId: *roomID,
			Name:   *name,
		},
	}
	client.Send(join)

	go func() {
		for {
			m, err := client.Recv()
			if err != nil {
				fmt.Printf("message recieved but could not be displayed: %v\n", err)
				return
			} else {
				fmt.Printf("%s: %s\n", m.GetName(), m.GetMessage())
			}
		}
	}()

	input := bufio.NewReader(os.Stdin)
	for {
		line, _, err := input.ReadLine()
		if err != nil {
			fmt.Printf("failed to read input: %v\n", err)
		}
		msg := &rpc.ChatMessage{
			Message: &rpc.Message{
				Message: string(line),
			},
		}
		if err := client.Send(msg); err != nil {
			fmt.Printf("failed to send message: %v\n", err)
		}
	}
}
