package main

import (
	"fmt"
	"time"
)

const size = 5 //number of philosophers and forks

var chans [size]chan int //channel
var eat [size]int
var stateFork [size]int

func philosopher(i int) {
	eat[i] = 0
	for eat[i] < 3 { //eat at least 3 times
		if (<-chans[i]) == -1 {
			//fmt.Println("I'm philosopher", i, "and I found the fork", i, "is free")
			chans[i] <- i //get the fork i
			//fmt.Println("I'm philosopher", i, "and I received the fork", i)
		} else {
			chans[i] <- -2 //-2 means do nothing
		}
		if (<-chans[(i+1)%size]) == -1 {
			//fmt.Println("I'm philosopher", i, "and I found the fork", i+1, "is free")
			chans[(i+1)%size] <- i //get the fork i+1
			//fmt.Println("I'm philosopher", i, "and I received the fork", (i + 1))
		} else {
			chans[(i+1)%size] <- -2 //-2 means do nothing
		}
		l := <-chans[i]          //check the state the fork i
		r := <-chans[(i+1)%size] //check the state the fork i+1
		if l == i && r == i {    //eat and put down the fork i and i+1
			//fmt.Println("I'm philosopher", i, "and I eaten and give up the fork", (i))
			//fmt.Println("I'm philosopher", i, "and I eaten and give up the fork", (i+1)%size)
			eat[i] = eat[i] + 1
			fmt.Println("I'm philosopher", i, "and I have eaten", eat[i])
			time.Sleep(2000000 * time.Nanosecond)
			chans[i] <- -1
			chans[(i+1)%size] <- -1 //???
		} else if l != i && r == i { //eat nothing and put down the fork i+1
			time.Sleep(2000000 * time.Nanosecond)
			chans[i] <- -1
			chans[(i+1)%size] <- -1
			//fmt.Println("I'm philosopher", i, "and I give up the fork", (i + 1))
		} else if l == i && r != i { //eat nothing and put down the fork i
			time.Sleep(2000000 * time.Nanosecond)
			chans[i] <- -1
			chans[(i+1)%size] <- -1
			//fmt.Println("I'm philosopher", i, "and I give up the fork", (i))
		} else { //eat nothing
			chans[i] <- -2
			chans[(i+1)%size] <- -2
		}

	}
}

func fork(i int) {
	stateFork[i] = -1 //on the table
	for {
		chans[i] <- stateFork[i]
		r := <-chans[i]
		//fmt.Println("attemp", i, r)
		if r != -2 {
			stateFork[i] = r
		}
		//fmt.Println("result", i, stateFork[i])
	}
}

func main() {
	for i := range chans {
		chans[i] = make(chan int)
	}
	for i := range chans {
		go fork(i)
	}
	for i := range chans {
		go philosopher(i)
	}
	for i := 0; i < 200; i++ {
		fmt.Println(i)
		time.Sleep(1000000 * time.Nanosecond)
	}

}
