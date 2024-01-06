package main

import (
	"log"
	"time"
)

var chans [2]chan packet //channel

type node struct {
	id    int
	seq   int
	state int // 0 = waiting, 1 = syn sent, 2 = syn received, 3 = established
}

type packet struct {
	ptype     int // 0 = syn, 1 = ack, 2 = synack, 3 = close
	source_id int
	seq       int
	ack       int
}

func (n *node) connect(nodeid int) {
	n.state = 1
	chans[nodeid] <- packet{0, n.id, n.seq, -1}
	log.Printf("node %d sent syn", n.id)

}

func (n *node) disconnect(nodeid int) {
	chans[nodeid] <- packet{3, n.id, -1, -1}
	n.state = 0
	log.Printf("node %d sent close", n.id)
}

func (n *node) respond(nodeid int, syn int) {
	n.state = 2
	chans[nodeid] <- packet{2, n.id, n.seq, syn + 1}
	log.Printf("node %d sent synack", n.id)
}

func (n *node) acknowledge(nodeid int, ack int) {
	n.state = 3
	chans[nodeid] <- packet{1, n.id, -1, ack + 1}
	log.Printf("node %d sent ack", n.id)
}
func (n *node) listen() {
	for {
		p := <-chans[n.id]
		if p.ptype == 0 && n.state == 0 {
			time.Sleep(100 * time.Millisecond)
			n.respond(p.source_id, p.seq)
			log.Printf("node %d received syn, state changed to %d", n.id, n.state)
		} else if p.ptype == 2 && n.state == 1 {
			time.Sleep(100 * time.Millisecond)
			n.acknowledge(p.source_id, p.seq)
			log.Printf("node %d received synack, state changed to %d", n.id, n.state)

		} else if p.ptype == 1 && n.state == 2 {
			time.Sleep(100 * time.Millisecond)
			n.state = 3
			log.Printf("node %d received ack, state changed to %d", n.id, n.state)

		} else if p.ptype == 3 {
			time.Sleep(100 * time.Millisecond)
			n.state = 0
			log.Printf("node %d received close, state changed to %d", n.id, n.state)

		}
	}
}

func main() {
	chans[0] = make(chan packet)
	chans[1] = make(chan packet)
	a := node{0, 0, 0}
	b := node{1, 0, 0}
	go a.listen()
	go b.listen()
	a.connect(1)
	time.Sleep(2 * time.Second)
	a.disconnect(1)
	time.Sleep(2 * time.Second)
	b.connect(0)
	time.Sleep(2 * time.Second)
	b.disconnect(0)
	for {
	}

}
