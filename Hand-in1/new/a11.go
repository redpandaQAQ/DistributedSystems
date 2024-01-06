package main

import (
	"fmt"
	"time"
)

const size = 5 //number of philosophers and forks

var chans [size * 2]chan int //channel
var eat [size]int
var stateFork [size]int

func philosopher(i int) {
	eat[i] = 0
	for eat[i] < 3 { //eat at least 3 times
		if (<-chans[2*i]) == -1 {
			//fmt.Println("I'm philosopher", i, "and I found the fork", i, "is free")
			chans[2*i] <- i //get the fork i
			//fmt.Println("I'm philosopher", i, "and I received the fork", i)
		} else {
			chans[2*i] <- -2 //-2 means do nothing
		}
		//time.Sleep(10000 * time.Nanosecond)
		if (<-chans[2*i+1]) == -1 {
			//fmt.Println("I'm philosopher", i, "and I found the fork", i+1, "is free")
			chans[2*i+1] <- i //get the fork i+1
			//fmt.Println("I'm philosopher", i, "and I received the fork", (i + 1))
		} else {
			chans[2*i+1] <- -2 //-2 means do nothing
		}
		l := <-chans[2*i]     //check the state the fork i
		r := <-chans[2*i+1]   //check the state the fork i+1
		if l == i && r == i { //eat and put down the fork i and i+1
			//fmt.Println("I'm philosopher", i, "and I eaten and give up the fork", (i))
			//fmt.Println("I'm philosopher", i, "and I eaten and give up the fork", (i+size))
			eat[i] = eat[i] + 1
			fmt.Println("I'm philosopher", i, "and I have eaten", eat[i])
			//time.Sleep(2000000 * time.Nanosecond)
			chans[2*i] <- -1
			chans[2*i+1] <- -1 //???
		} else if l != i && r == i { //eat nothing and put down the fork i+1
			//time.Sleep(2000000 * time.Nanosecond)
			chans[2*i] <- -2
			chans[2*i+1] <- -1
			//fmt.Println("I'm philosopher", i, "and I give up the fork", (i + 1))
		} else if l == i && r != i { //eat nothing and put down the fork i
			//time.Sleep(2000000 * time.Nanosecond)
			chans[2*i] <- -1
			chans[2*i+1] <- -2
			//fmt.Println("I'm philosopher", i, "and I give up the fork", (i))
		} else { //eat nothing
			chans[2*i] <- -2
			chans[2*i+1] <- -2
		}

	}
}

func fork(i int) {
	stateFork[i] = -1 //on the table
	for {
		select {
		case chans[2*i] <- stateFork[i]:
			r := <-chans[2*i]
			if r != -2 {
				stateFork[i] = r
			}
		case chans[(2*i+9)%(size*2)] <- stateFork[i]:
			r := <-chans[(2*i+9)%(size*2)]
			if r != -2 {
				stateFork[i] = r
			}

		}
	}
}

func main() {
	for i := range chans {
		chans[i] = make(chan int)
	}
	for i := 0; i < size; i++ {
		go fork(i)
	}
	for i := 0; i < size; i++ {
		go philosopher(i)
	}
	for i := 0; i < 200; i++ {
		fmt.Println(i)
		time.Sleep(100000 * time.Nanosecond)
	}

}
