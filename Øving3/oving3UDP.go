package main

import (
    "fmt"
    "net"
    "os"
		"time"
)

/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func Read(connection *net.UDPConn){
	buf := make([]byte, 1024)

	for {
			n,addr,err := connection.ReadFromUDP(buf)
			fmt.Println("Received ",string(buf[0:n]), " from ",addr)

			if err != nil {
					fmt.Println("Error: ",err)
			}
	}
}

func Send(connection *net.UDPConn){

	for {
		fmt.Println("msg sendt")
		connection.Write([]byte("Hello from client21"))
		time.Sleep(time.Second * 1)
	}
}

func main() {

    ServerAddr,err := net.ResolveUDPAddr("udp","129.241.187.43:20021")
    CheckError(err)
		ReadPort, err := net.ResolveUDPAddr("udp",":20021")
		CheckError(err)


    ClientSend, err := net.DialUDP("udp", nil, ServerAddr)
    CheckError(err)

		ClientRead, err := net.ListenUDP("udp", ReadPort)
    CheckError(err)

		go Read(ClientRead)
		go Send(ClientSend)

		deadChan :=make(chan bool, 1)
		<- deadChan


}
