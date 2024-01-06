package main

import (
	proto "chittyChat/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"google.golang.org/grpc"
)

// Struct that will be used to represent the Server.
type Server struct {
	proto.UnimplementedChatServer // Necessary
	name                          string
	port                          int
	clock                         Clock
	messageChannel                map[int](chan *proto.ChatMessage) //map from client id to channel
}

type Clock struct {
	clock int
	mu    sync.Mutex
}

// Used to get the user-defined port for the server from the command line
//var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	//log.SetOutput(os.Stdout)
	//flag.Parse()
	// Create a server struct
	server := &Server{
		name:           "serverName",
		port:           int(5454),
		clock:          Clock{clock: 0},
		messageChannel: make(map[int](chan *proto.ChatMessage)),
	}
	for i, _ := range server.messageChannel {
		server.messageChannel[i] = make(chan *proto.ChatMessage, 100) //buffer size 100
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterChatServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (c *Server) ClientSend(ctx context.Context, in *proto.ChatMessage) (*proto.ServerResponse, error) {
	log.Printf("Server received a message from Client %d\n", in.SenderId)
	c.UpdateClock(in.Timestamp)
	c.ClockIncrease() //increase clock when receive message from client
	for i := range c.messageChannel {
		go c.Broadcast(i, int(in.SenderId), in)
	}
	return &proto.ServerResponse{
		Timestamp: int64(c.GetClock()),
		Result:    true,
	}, nil
}

func (c *Server) ClientReceive(in *proto.ClientRequest, stream proto.Chat_ClientReceiveServer) error {
	//log.Printf("Client %d request new message\n", in.ClientId)
	for {
		message := <-c.messageChannel[int(in.ClientId)]
		//c.ClockIncrease() //increase clock when send message to client
		message.Timestamp = int64(c.GetClock())
		if err := stream.Send(message); err != nil {
			return err
		}
	}
	//return nil
}

func (c *Server) Join(ctx context.Context, in *proto.ClientRequest) (*proto.ServerResponse, error) {
	//log.Printf("Client %d connected the chat\n", in.ClientId)
	c.messageChannel[int(in.ClientId)] = make(chan *proto.ChatMessage)
	c.UpdateClock(in.Timestamp)
	c.ClockIncrease()
	text := fmt.Sprintf("Participant %d joined Chitty-Chat at Lamport time %d\n", in.ClientId, c.GetClock())
	message := &proto.ChatMessage{
		SenderId:    int64(in.ClientId),
		Timestamp:   int64(c.GetClock()),
		MessageText: text,
	}
	for i := range c.messageChannel {
		go c.Broadcast(i, int(in.ClientId), message)
	}
	/*
		for i, mc := range c.messageChannel {
			if i != int(in.ClientId) {
				text := fmt.Sprintf("Participant %d joined Chitty-Chat at Lamport time %d\n", in.ClientId, c.GetClock())
				mc <- &proto.ChatMessage{
					SenderId:    int64(in.ClientId),
					Timestamp:   int64(c.GetClock()),
					MessageText: text,
				}
			}
		}
	*/
	return &proto.ServerResponse{
		Timestamp: int64(c.GetClock()),
		Result:    true,
	}, nil
}

func (c *Server) Leave(ctx context.Context, in *proto.ClientRequest) (*proto.ServerResponse, error) {
	//log.Printf("Client %d disconnected to the server\n", in.ClientId)
	c.UpdateClock(in.Timestamp)
	c.ClockIncrease()
	for i, mc := range c.messageChannel {
		if i != int(in.ClientId) {
			text := fmt.Sprintf("Participant %d left Chitty-Chat at Lamport time %d", in.ClientId, c.GetClock())
			mc <- &proto.ChatMessage{
				SenderId:    int64(in.ClientId),
				Timestamp:   int64(c.GetClock()),
				MessageText: text,
			}
		}
	}
	delete(c.messageChannel, int(in.ClientId))
	return &proto.ServerResponse{
		Timestamp: int64(c.GetClock()),
		Result:    true,
	}, nil
}
func (c *Server) UpdateClock(ts int64) {
	c.clock.mu.Lock()
	defer c.clock.mu.Unlock()
	if c.clock.clock < int(ts) {
		c.clock.clock = int(ts)
	}
	c.clock.clock++
	log.Printf("Server clock is updated to %d\n", c.clock.clock)
}

func (c *Server) ClockIncrease() {
	c.clock.mu.Lock()
	defer c.clock.mu.Unlock()
	c.clock.clock++
	log.Printf("Server clock is updated to %d\n", c.clock.clock)
}

func (c *Server) Broadcast(i int, sender int, in *proto.ChatMessage) {
	//if i != sender {
	c.messageChannel[i] <- in
	//}
}
func (c *Server) GetClock() int {
	c.clock.mu.Lock()
	defer c.clock.mu.Unlock()
	return c.clock.clock
}
