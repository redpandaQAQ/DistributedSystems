package main

import (
	"fmt"
	"time"
)

var clientToServer = make(chan Packet)
var serverToClient = make(chan Packet)

func client(i int) {
	//defer wg.Done()

	state := StateSynSent
	seq := 100 * i
	// Client FSM
	for {
		time.Sleep(1 * time.Second)
		switch state {
		case StateSynSent:
			// Client sending SYN packet
			seq++
			clientToServer <- Packet{Type: "syn", Seq: seq}
			state = StateSynReceived

		case StateSynReceived:
			// Client receiving SYN+ACK packet
			p := <-serverToClient
			fmt.Printf("Client got: syn ack=%d seq=%d\n", p.Ack, p.Seq)
			state = StateEstablished

			// Client sending ACK packet
			seq++
			clientToServer <- Packet{Type: "ack", Ack: p.Seq + 1, Seq: seq}

		case StateEstablished:
			// Simulate data transfer
			p := <-serverToClient
			fmt.Printf("Client got: seq=%d data=%s\n", p.Seq, p.Data)
			state = StateSynSent
		}
	}
}

// Packet represents a TCP packet.
type Packet struct {
	Type string
	Syn  int
	Ack  int
	Seq  int
	Data string
}

// ConnectionState represents the state of the connection FSM.
type ConnectionState int

const (
	// StateSynSent represents the SYN sent state.
	StateSynSent ConnectionState = iota
	// StateSynReceived represents the SYN received state.
	StateSynReceived
	// StateEstablished represents the established state.
	StateEstablished
)

func main() {
	// Create channels of type Packet

	// Wait group to synchronize goroutines
	//var wg sync.WaitGroup

	// Client FSM goroutine
	//wg.Add(1)
	go client(1)
	//go client(2)

	// Server FSM goroutine
	//wg.Add(1)
	go func() {
		//defer wg.Done()

		state := StateSynSent
		seq := 300
		// Server FSM
		for {
			time.Sleep(1 * time.Second)
			switch state {
			case StateSynSent:
				// Server receiving SYN packet
				p := <-clientToServer
				fmt.Printf("Server got: syn seq=%d\n", p.Seq)
				state = StateSynReceived

				// Server sending SYN+ACK packet
				seq++
				serverToClient <- Packet{Type: "syn", Ack: p.Seq + 1, Seq: seq}

			case StateSynReceived:
				// Server receiving ACK packet
				p := <-clientToServer
				fmt.Printf("Server got: ack=%d seq=%d\n", p.Ack, p.Seq)
				state = StateEstablished
				seq++
				// Simulate data transfer
				serverToClient <- Packet{Type: "data", Seq: seq, Data: "Hello from server"}

			case StateEstablished: // just exit, maybe can enter listening (idle) status
				state = StateSynSent
			}
		}
	}()

	// Wait for both goroutines to finish
	//wg.Wait()
	for {

	}
}

// func main() {
// 	// Create channels of type Packet
// 	clientToServer := make(chan Packet)
// 	serverToClient := make(chan Packet)

// 	//  wait for a fixed number of goroutines to complete their work
// 	var wg sync.WaitGroup

// 	// Client goroutine
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()

// 		// Client sending SYN packet
// 		clientToServer <- Packet{Type: "syn", Seq: 1}

// 		// Client receiving SYN+ACK packet
// 		p := <-serverToClient
// 		fmt.Printf("Client got: syn ack=%d seq=%d\n", p.Ack, p.Seq)

// 		// Client sending ACK packet
// 		clientToServer <- Packet{Type: "ack", Ack: p.Seq + 1, Seq: 2}

// 		// Simulate data transfer
// 		p = <-serverToClient
// 		clientToServer <- Packet{Type: "data", Seq: 3, Data: "Hello from client"}
// 		fmt.Printf("Client got: seq=%d data=%s\n", p.Seq, p.Data)
// 	}()

// 	// Server goroutine
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()

// 		// Server receiving SYN packet
// 		p := <-clientToServer
// 		fmt.Printf("Server got: syn seq=%d\n", p.Seq)

// 		// Server sending SYN+ACK packet
// 		serverToClient <- Packet{Type: "syn", Ack: p.Seq + 1, Seq: 1}

// 		// Server receiving ACK packet
// 		p = <-clientToServer
// 		fmt.Printf("Server got: ack=%d seq=%d\n", p.Ack, p.Seq)

// 		// Simulate data transfer
// 		serverToClient <- Packet{Type: "data", Seq: 2, Data: "Hello from server"}
// 		p = <-clientToServer
// 		fmt.Printf("Server got: seq=%d data=%s\n", p.Seq, p.Data)
// 	}()

// 	// Wait for both goroutines to finish
// 	wg.Wait()
// }
