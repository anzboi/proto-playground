package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/anzboi/proto-playground/pkg/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServiceImpl struct {
	// map from room ID to chat room
	rooms map[string]*Room
}

func NewChatService() *ChatServiceImpl {
	return &ChatServiceImpl{
		rooms: map[string]*Room{},
	}
}

// Room defines a single chat room
type Room struct {
	// record of all messages
	messages []*rpc.Message

	// list of named subscribers
	subscribers map[string]chan<- *rpc.Message

	closed bool
}

func NewRoom() *Room {
	return &Room{
		messages:    []*rpc.Message{},
		subscribers: map[string]chan<- *rpc.Message{},
	}
}

// Push a new message onto the chat room
func (t *Room) Push(pusher string, m *rpc.Message) {
	if t == nil {
		return
	}
	t.messages = append(t.messages, m)
	for name, sub := range t.subscribers {
		// Don't send the message back to the sender
		if name != pusher {
			sub <- m
		}
	}
}

// Subscribe to a chat room
//
// must provide a chanel that will be used to stream new messages as they arrive
func (t *Room) Subscribe(name string, ch chan<- *rpc.Message) {
	t.subscribers[name] = ch
	joinedMessage := &rpc.Message{
		Name:    name,
		Message: "Has joined the chat",
	}
	t.Push(name, joinedMessage)
}

// Unsubscribe removes a channel from a chat room
//
// MUST unsubscribe BEFORE closing a channel
func (t *Room) Unsubscribe(name string) {
	if t == nil {
		return
	}

	ch, _ := t.subscribers[name]
	delete(t.subscribers, name)
	close(ch)
	leftMessage := &rpc.Message{
		Name:    name,
		Message: "Has left the chat",
	}
	t.Push(name, leftMessage)
}

func (t *Room) Close() {
	for _, ch := range t.subscribers {
		closingMessage := &rpc.Message{
			Name:    "chat server",
			Message: "chat room is closing",
		}
		ch <- closingMessage
		close(ch)
	}
	t.closed = true
	t.messages = nil
	t.subscribers = nil
}

// Chat implements the bidirectional chat rpc
func (i *ChatServiceImpl) Chat(stream rpc.ChatService_ChatServer) error {
	// Recieve the join parameters
	first, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, "Failed to read initial join request")
	}

	// Join a chat room
	username := first.GetJoin().GetName()
	roomID := first.GetJoin().GetRoomId()
	var room *Room
	if Room, ok := i.rooms[roomID]; ok {
		room = Room
	} else {
		return status.Errorf(codes.NotFound, "room %s does not exist", roomID)
	}
	ch := make(chan *rpc.Message)
	room.Subscribe(username, ch)
	log.Printf("%s joined chat room %s", username, roomID)

	// Stream previous 10 messages
	prev := make([]*rpc.Message, 0, 10)
	for i := len(room.messages) - 1; i >= 0 && i >= len(room.messages)-10; i-- {
		prev = append(prev, room.messages[i])
	}
	for i := len(prev) - 1; i >= 0; i-- {
		stream.Send(prev[i])
	}

	// Listen for other chatters on the Room
	closed := make(chan struct{})
	go func() {
		for m := range ch {
			_ = stream.Send(m)
		}
		closed <- struct{}{}
	}()

	go func() {
		for {
			m, err := stream.Recv()
			if err != nil {
				return
			}
			if room.closed {
				return
			}
			if message := m.GetMessage(); message != nil {
				message.Name = username // make sure the messages name matches the username
				room.Push(username, message)
			}
		}
	}()

	select {
	case <-closed:
		return status.Error(codes.OK, "chat room closed")
	case <-stream.Context().Done():
		room.Unsubscribe(username)
		log.Printf("%s left chat room %s", username, roomID)
		return nil
	}
}

func (i *ChatServiceImpl) CreateChatRoom(ctx context.Context, req *rpc.CreateChatRoomRequest) (*rpc.ChatRoom, error) {
	roomID := newRoomID()
	i.rooms[roomID] = NewRoom()
	log.Printf("created chat room %s", roomID)
	return &rpc.ChatRoom{RoomId: roomID}, nil
}

func (i *ChatServiceImpl) ListChatRooms(ctx context.Context, req *rpc.ListChatRoomsRequest) (*rpc.ListChatRoomsResponse, error) {
	resp := make([]*rpc.ChatRoom, 0, len(i.rooms))
	for roomID, _ := range i.rooms {
		resp = append(resp, &rpc.ChatRoom{RoomId: roomID})
	}
	return &rpc.ListChatRoomsResponse{ChatRooms: resp}, nil
}

func (i *ChatServiceImpl) DeleteChatRoom(ctx context.Context, req *rpc.DeleteChatRoomRequest) (*rpc.ChatRoom, error) {
	room, ok := i.rooms[req.GetRoomId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat room %s does not exist", req.GetRoomId())
	}
	room.Close()
	delete(i.rooms, req.GetRoomId())
	log.Printf("chat room %s closed", req.GetRoomId())
	return &rpc.ChatRoom{RoomId: req.GetRoomId()}, nil
}

func newRoomID() string {
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
