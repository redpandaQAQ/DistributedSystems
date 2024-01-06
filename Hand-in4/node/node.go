package main

import (
	"bufio"
	"context"
	"flag"
	proto "handin4/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Struct that will be used to represent the Server.
type Node struct {
	proto.UnimplementedMutualExclusionServer // Necessary
	id                                       int
	port                                     int
	state                                    bool
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")            //port of current node
var id = flag.Int("id", 0, "server id")                         //id of current node
var nextPort = flag.Int("nextport", 0, "next node port number") //port of next node
var nextId = flag.Int("nextid", 0, "next node id")              //id of next node

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	node := &Node{
		id:    *id,
		port:  *port,
		state: false,
	}

	// Start the server
	go startNode(node) // Start the server in a goroutine so it does not block.
	if *id == 1 {      //Only node with id=1 can start sending token
		go SendToken()
	}
	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("type x to request access to the Critical Section")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "X" || input == "x" {
			node.state = true //Set state to true to indicate that node is requesting the critical section
			log.Printf("Node %d is requesting access to the critical section", node.id)

		} else {
			log.Print("Unknown command, ignored")
		}

	}
}

func startNode(node *Node) {
	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(node.port))
	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started node %d at port: %d\n", node.id, node.port)

	// Register the grpc server and serve its listener
	proto.RegisterMutualExclusionServer(grpcServer, node)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not create serve listener")
	}
}

func connectToServer() (proto.MutualExclusionClient, error) {
	// Dial the ndoe at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*nextPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to node %d port %d", *nextId, *nextPort)
	} else {
		//log.Printf("Connected to the node %d\n", *nextPort)
	}
	return proto.NewMutualExclusionClient(conn), nil
}
func SendToken() {
	serverConnection, _ := connectToServer()
	Response, err := serverConnection.SendToken(context.Background(), &proto.Token{
		Id: int32(*id),
	})
	if err != nil {
		log.Printf("Could not sent token to node %d, will retry in 1 second", *nextId) //If node is not up, retry in 1 second
		time.Sleep(time.Second)
		SendToken()
	} else if Response.Success {
		log.Printf("Token was sent to %d", *nextId) //don't know why this is not working
	}
}

func (n *Node) SendToken(ctx context.Context, in *proto.Token) (*proto.Response, error) {
	if n.state {
		log.Printf("Node %d is now in the critical section", n.id) //If node is requesting the critical section, enter it
		time.Sleep(5 * time.Second)                                //Wait 5 seconds to simulate work
		n.state = false                                            //Exit critical section
		log.Printf("Node %d exits critical section", n.id)
		SendToken() //Always send token to next node
		return &proto.Response{
			Success: true,
		}, nil
	} else {
		log.Printf("Node %d is not requesting the critical section, give token to node %d", n.id, *nextId)
		time.Sleep(5 * time.Second) //Wait 5 seconds before sending token
		SendToken()                 //Always send token to next node
		return &proto.Response{
			Success: true,
		}, nil
	}
}
