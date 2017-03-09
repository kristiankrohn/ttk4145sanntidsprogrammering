package main

import (
	. "./Kontrollmodul/Heismodul"
	//. "./Kontrollmodul/Heismodul/driver"
	. "./Kontrollmodul/Nettverksmodul"
	. "./Kontrollmodul"
	"fmt"
	"runtime"
	//"strconv"
	//"strings"
	"time"
	//"os"
	//"encoding/gob"
	"net"
	"os/exec"
	"encoding/binary"
)

func Watchdog()(uint64){
	var counter uint64 = 0
	buffer := make([]byte, 8)
	ListenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.255:20221")
	CheckError(err)
	Listener, err := net.ListenUDP("udp", ListenAddr)
	CheckError(err)

	fmt.Println("Connection established")
	fmt.Println("A new elevatorcontoller has been born")
	for {
		Listener.SetReadDeadline(time.Now().Add(time.Second*2))
		n, _, err := Listener.ReadFromUDP(buffer)

		if err != nil{
			break
		} else{
			counter = binary.BigEndian.Uint64(buffer[0:n])
		}
	}
	Listener.Close()

	command := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err = command.Run()
	CheckError(err)
	fmt.Println("My father is dead, i'm now in control")

	return counter

}

func IsAlive(counter uint64){
	isAliveAddr, err := net.ResolveUDPAddr("udp", "127.0.0.255:20221")
	CheckError(err)
	isAlive, err := net.DialUDP("udp", nil, isAliveAddr)
	CheckError(err)
	for{
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, counter)
		_,_ = isAlive.Write(buffer)
		time.Sleep(time.Second * 1)
	}
}

func main() {

	var counter uint64 = Watchdog()
	

	nextFloor := make(chan int, 20)
	orderFinished := make(chan bool, 5)
	up_button := make(chan int, 4)
	down_button := make(chan int, 4)
	internal_button := make(chan int, 4)
	message := make(chan string, 20)
	recievedmessage := make(chan string, 40)

	runtime.GOMAXPROCS(runtime.NumCPU())
	go IsAlive(counter)
	Init_system(nextFloor)
	fmt.Println("Init finished")

	go Broadcast(message, recievedmessage)
	go Elevator_driver(nextFloor, orderFinished)
	go TCP_sender(message, recievedmessage)
	go Local_orders(internal_button, nextFloor, orderFinished)
	go Handle_buttons(up_button, down_button, internal_button)
	go External_orders(message, up_button, down_button, nextFloor)
	go Displayfloor()
	go TCP_listener(recievedmessage)
	go Incomming_message(recievedmessage, message)
	go Assess_cost(nextFloor)
	go Clear_orders(orderFinished, nextFloor, message)
	go Resend_externalorders(message)

	
	

	deadChan := make(chan bool, 1)
	<-deadChan
}
