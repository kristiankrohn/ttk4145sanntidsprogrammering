package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
	//"bytes"
)

const numberofelevators int = 255

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func Broadcast() {

	BroadcastAddr, err := net.ResolveUDPAddr("udp", "129.241.187.255:20021")
	CheckError(err)
	BroadcastSpammer, err := net.DialUDP("udp", nil, BroadcastAddr)
	CheckError(err)

	for {
		BroadcastSpammer.Write([]byte("Hello, i'm an elevator from group 67 - connect to me!"))
		time.Sleep(time.Second * 1)
	}
}

func TCP_sender(message chan string) {

	var numberofIPs int = 0
	var addressArray [numberofelevators]string

	ListenPort, err := net.ResolveUDPAddr("udp", ":20021")
	CheckError(err)
	Listen, err := net.ListenUDP("udp", ListenPort)
	CheckError(err)
	buf := make([]byte, 1024)

	for {

		n, addr, err := Listen.ReadFromUDP(buf)
		//fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
			Listen.Close()
			net.ListenUDP("udp", ListenPort)

		}

		firsthalf := strings.Split(addr.String(), ":")
		NewIP := firsthalf[0]
		var IPmatch bool = false

		if string(buf[0:n]) == "Hello, i'm an elevator from group 67 - connect to me!" {

			for i := 0; i <= numberofIPs; i++ {
				if addressArray[i] == NewIP {
					IPmatch = true
				}
			}

			if IPmatch == false {
				numberofIPs++
				addressArray[numberofIPs] = NewIP

				fmt.Println("New machine at: ", addressArray[numberofIPs])

			}
		}

		select {
		case sendmessage := <-message:
			{
				for i := 0; i <= numberofIPs; i++ {
					Clientaddress, err := net.ResolveTCPAddr("tcp", string(addressArray[i]+":20021"))
					CheckError(err)
					Client, err := net.DialTCP("tcp", nil, Clientaddress)
					if err != nil {
						fmt.Println("Disconnect :", addressArray[i])
						for j := 0; j <= (numberofIPs - i); j++ {
							addressArray[i+j] = addressArray[i+j+1]
						}
						numberofIPs--

					} else {

						_, err = Client.Write([]byte(sendmessage))
						CheckError(err)
					}

				}
			}
		default:
		}
	}

}

func TCP_listener() {
	listenPort, err := net.Listen("tcp", ":20021")
	CheckError(err)

	for {
		connection, err := listenPort.Accept()
		CheckError(err)
		addr := connection.RemoteAddr()

		firsthalf := strings.Split(addr.String(), ":")
		IP := firsthalf[0]

		if IP != "127.0.0.1" {
			buf := make([]byte, 1024)
			n, err := connection.Read(buf)
			CheckError(err)

			address := connection.RemoteAddr()
			fmt.Println("Recieved message : ", string(buf[0:n]), " from ", address)
			connection.Close()
		} else {
			connection.Close()
		}
	}

}

func Test(message chan string) {

	for {
		message <- string("Hello")
		time.Sleep(time.Second * 1)
	}
}

func main() {

	message := make(chan string, 1024)
	runtime.GOMAXPROCS(runtime.NumCPU())

	go Broadcast()
	go TCP_sender(message)
	go TCP_listener()
	go Test(message)

	deadChan := make(chan bool, 1)
	<-deadChan
}
