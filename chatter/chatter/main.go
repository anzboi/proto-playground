package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/anzboi/proto-playground/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	create   = flag.Bool("create", false, "create a chat room")
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

	chatService := rpc.NewChatServiceClient(cc)

	if *create {
		resp, err := chatService.CreateChatRoom(context.Background(), &rpc.CreateChatRoomRequest{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Chat room created: %s\n", resp.GetRoomId())
		return
	}

	rooms, err := chatService.ListChatRooms(context.Background(), &rpc.ListChatRoomsRequest{})
	if err != nil {
		panic(err)
	}

	input := bufio.NewReader(os.Stdin)
	fmt.Println("Chat rooms:")
	roomMap := map[int]string{}
	for i, room := range rooms.GetChatRooms() {
		roomMap[i+1] = room.GetRoomId()
		fmt.Printf("[%d]: %s\n", i+1, room.GetRoomId())
	}
	fmt.Print("Pick a room: ")
	line, _, err := input.ReadLine()
	if err != nil {
		panic(err)
	}
	r, err := strconv.Atoi(string(line))
	if err != nil {
		fmt.Println("error: must input a number")
		return
	}
	roomID, ok := roomMap[r]
	if !ok {
		fmt.Println("error: must select a room that exists")
		return
	}

	fmt.Print("Enter Name tag: ")
	nameLine, _, err := input.ReadLine()
	if err != nil {
		panic(err)
	}
	name := string(nameLine)

	client, err := chatService.Chat(context.Background())
	if err != nil {
		panic(err)
	}

	join := &rpc.ChatMessage{
		Join: &rpc.JoinChat{
			RoomId: roomID,
			Name:   name,
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
