package main

import (
	. "./Controlmodule"
	. "./Controlmodule/Elevatormodule"
	. "./Controlmodule/Networkmodule"
	"encoding/binary"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

func Watchdog() uint64 {
	var counter uint64 = 0
	buffer := make([]byte, 8)
	ListenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.255:20221")
	CheckError(err)
	Listener, err := net.ListenUDP("udp", ListenAddr)
	CheckError(err)

	fmt.Println("Connection established")
	fmt.Println("A new elevatorcontoller has been born")

	isAliveAddr, err := net.ResolveUDPAddr("udp", "127.0.0.255:22221")
	CheckError(err)
	isAlive, err := net.DialUDP("udp", nil, isAliveAddr)
	CheckError(err)



	for {

		binary.BigEndian.PutUint64(buffer, counter)
		_, _ = isAlive.Write(buffer)

		Listener.SetReadDeadline(time.Now().Add(time.Second * 1))
		n, _, err := Listener.ReadFromUDP(buffer)
		if err != nil {
			break
		} else {
			counter = binary.BigEndian.Uint64(buffer[0:n])
			time.Sleep(time.Millisecond * 10)
		}

	}
	Listener.Close()
	isAlive.Close()

	command := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
	err = command.Run()
	CheckError(err)
	fmt.Println("I'm now in control")

	return counter

}

func Watchcat(counter uint64) {
	isAliveAddr, err := net.ResolveUDPAddr("udp", "127.0.0.255:20221")
	CheckError(err)
	isAlive, err := net.DialUDP("udp", nil, isAliveAddr)
	CheckError(err)

	for {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, counter)
		_, _ = isAlive.Write(buffer)

		time.Sleep(time.Millisecond * 333)
	}
}

func IsDogAlive(){
	ListenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.255:22221")
	CheckError(err)
	Listener, err := net.ListenUDP("udp", ListenAddr)
	CheckError(err)
	buffer := make([]byte, 8)
	for{
		Listener.SetReadDeadline(time.Now().Add(time.Second * 2))
		_, _, err := Listener.ReadFromUDP(buffer)
		if err != nil {
			command := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run main.go")
			err = command.Run()
			CheckError(err)
			fmt.Println("Watchdog has died")

		}
		time.Sleep(time.Millisecond * 10)
	}
}

func main() {
	nextFloor := make(chan int, 20)
	go Display_floor()

	var counter uint64 = Watchdog()
	go Watchcat(counter)
	go IsDogAlive()
	Init_system(nextFloor)

	orderFinished := make(chan bool, 5)
	upButton := make(chan int, 4)
	downButton := make(chan int, 4)
	internalButton := make(chan int, 4)
	message := make(chan string, 20)
	recievedMessage := make(chan string, 40)
	stopElevator := make(chan bool, 5)

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Init finished")

	go Broadcast(message, recievedMessage)
	go Elevator_driver(nextFloor, orderFinished, stopElevator)
	go TCP_sender(message, recievedMessage)
	go Local_orders(internalButton, nextFloor, orderFinished)
	go Handle_buttons(upButton, downButton, internalButton)
	go External_orders(message, upButton, downButton, nextFloor)
	go TCP_listener(recievedMessage)
	go Message_handler(recievedMessage, message)
	go Assess_cost(nextFloor)
	go Clear_orders(orderFinished, nextFloor, message, stopElevator)
	go Resend_externalorders(message)

	deadChan := make(chan bool, 1)
	<-deadChan
}
