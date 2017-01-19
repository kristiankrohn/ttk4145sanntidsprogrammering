// Go 1.2
// go run helloworld_go.go

package main

import (
    ."fmt"
    ."runtime"
    //."time"
)

var global_i int = 0

func thread_1(ch chan int, doneChannel chan bool) {
	for j := 0; j < 1000; j++ {
    i := <- ch
		i++
    Println("thread_1")
    ch <- i
	}
  doneChannel <- true
}

func thread_2(ch chan int, doneChannel chan bool) {
	for k := 0; k < 1000; k++ {
    i := <- ch
		i--
    Println("thread_2")
    ch <- i
	}
  doneChannel <- true
}

func main() {
  //GOMAXPROCS(2)

  channel := make(chan int, 1)
  doneChannel := make(chan bool, 1)

  GOMAXPROCS(2)
  go thread_1(channel, doneChannel)           // This spawns thread_1() as a goroutine
	go thread_2(channel, doneChannel)					  // This spawns thread_2() as a goroutine
  channel <- global_i
  <- doneChannel
  <- doneChannel
  //Sleep(100*Millisecond)
  Println(<-channel)
}
