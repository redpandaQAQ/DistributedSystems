package main

import (
	"bufio"
	proto "chittyChat/grpc"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id    int
	clock Clock
}

type Clock struct {
	clock int
	mu    sync.Mutex
}

var (
	//serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
	cid            = flag.Int("cid", 0, "clientid") // client id
	serverPort int = int(5454)
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()
	// Create a client
	client := &Client{
		id:    *cid,
		clock: Clock{clock: 0},
	}
	joinChat(client)
	go receiveMessage(client)
	sendMessage(client)
	leaveChat(client)
}

func sendMessage(client *Client) {
	// Connect to the server
	serverConnection, _ := connectToServer()

	// Wait for input in the client terminal
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputext := scanner.Text()
		if inputext == "leave" {
			break
		}
		client.ClockPlusOne()
		//log.Printf("Client send message: %s\n", inputext)
		text := fmt.Sprintf("From %d at clock %d: %s", client.id, client.GetClock(), inputext)
		_, err := serverConnection.ClientSend(context.Background(), &proto.ChatMessage{
			SenderId:    int64(client.id),
			Timestamp:   int64(client.GetClock()),
			MessageText: text,
		})

		if err != nil {
			log.Print(err.Error())
			client.ClockMinusOne()
		} else {
			//log.Printf("Server received")
		}

	}
}

func connectToServer() (proto.ChatClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", serverPort)
	} else {
		//log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewChatClient(conn), nil
}

func receiveMessage(client *Client) {
	// Connect to the server
	serverConnection, _ := connectToServer()
	stream, err := serverConnection.ClientReceive(context.Background(), &proto.ClientRequest{
		ClientId:  int64(client.id),
		Timestamp: int64(client.GetClock()),
	})
	if err != nil {
		log.Print(err.Error())
	}
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("failed: %v", err)
		}
		client.UpdateClock(message.Timestamp)
		log.Printf("%s, current clock is %d", message.MessageText, client.GetClock())
	}

}

func joinChat(client *Client) {
	serverConnection, _ := connectToServer()
	client.ClockPlusOne()
	returnMessage, err := serverConnection.Join(context.Background(), &proto.ClientRequest{
		ClientId:  int64(client.id),
		Timestamp: int64(client.GetClock()),
	})
	if err != nil {
		log.Print(err.Error())
	} else {
		if returnMessage.Result {
			log.Printf("Client %d join chat\n", client.id)
		} else {
			log.Print(err.Error())
		}

	}
}

func leaveChat(client *Client) {
	serverConnection, _ := connectToServer()
	client.ClockPlusOne()
	returnMessage, err := serverConnection.Leave(context.Background(), &proto.ClientRequest{
		ClientId:  int64(client.id),
		Timestamp: int64(client.GetClock()),
	})
	if err != nil {
		log.Print(err.Error())
	} else {
		if returnMessage.Result {
			log.Printf("Client %d left chat\n", client.id)
		} else {
			log.Print(err.Error())
		}

	}
}

func (c *Client) UpdateClock(ts int64) {
	c.clock.mu.Lock()
	defer c.clock.mu.Unlock()
	if c.clock.clock < int(ts) {
		c.clock.clock = int(ts)
	}
	c.clock.clock++
	log.Printf("Client%d clock is updated to %d\n", c.id, c.clock.clock)
}

func (c *Client) ClockPlusOne() {
	c.clock.mu.Lock()
	defer c.clock.mu.Unlock()
	c.clock.clock++
	log.Printf("Client%d clock is updated to %d\n", c.id, c.clock.clock)
}
func (c *Client) ClockMinusOne() {
	c.clock.mu.Lock()
	defer c.clock.mu.Unlock()
	c.clock.clock--
	log.Printf("Client%d clock is updated to %d\n", c.id, c.clock.clock)
}
func (c *Client) GetClock() int {
	c.clock.mu.Lock()
	defer c.clock.mu.Unlock()
	return c.clock.clock
}
