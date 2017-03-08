package Nettverksmodul

import (
	"fmt"
	"net"
	//"os"
	//"runtime"
	"strings"
	"time"
	//"bytes"
	"strconv"


)

/* 			 Forslag til meldingsoppbygging
	//	eksempel:
	//	new order at floor 3		:		message(0,3,0,ipAddr)
	//	cost				:		message(1,0,24,ipAddr)
	//	completeOrder at floor 3	:		message(2,3,0,ipAddr)

type message struct {
	messageType 	int
	floor		int
	cost		int
	ipAddr		int
}

const (
	newOrder	int = iota
	cost
	completeOrder
)

*/

const Numberofelevators int = 10 


func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		//os.Exit(0)
	}
}

func Last_byte_of_my_IP() int{ //Borrowed from https://github.com/TTK4145/Network-go/blob/master/network/localip/localip.go
	var localIP string
	connAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		err = nil
		return 0
	}	
	conn, err := net.DialUDP("udp", nil, connAddr)
	if err != nil {
		err = nil
		return 0
	}	
	localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	lastaddressbytestring := strings.Split(localIP, ".")
	conn.Close()
	lastaddressbyte_i64, err:= strconv.ParseInt(lastaddressbytestring[3], 10, 64)
	CheckError(err)
	lastaddressbyte := int(lastaddressbyte_i64)
	return lastaddressbyte
}

func Broadcast(message chan string, recievedmessage chan string) {
	// rewrite to statemachine to handle reconnects
	for{
			BroadcastAddr, err := net.ResolveUDPAddr("udp", "129.241.187.255:20021")
			if err != nil {
				fmt.Println("Warning: ", err)
			}	
			BroadcastSpammer, err := net.DialUDP("udp", nil, BroadcastAddr)
			if err != nil{
				fmt.Println("Warning: ", err)
				buffermessage := <- message // message loopback, remove when TCP sender throws an no connection erre
				
				recievedmessage <- strings.Join([]string{buffermessage, "0.0.0.255:9000"}, ",")
			} else {
					BroadcastSpammer.Write([]byte("Hello, i'm an elevator from group 67 - connect to me!"))
					BroadcastSpammer.Close()
			}
		time.Sleep(time.Second * 1)
	}
}

func TCP_sender(message chan string, recievedmessage chan string) {

	var numberofIPs int = 0
	var addressArray [Numberofelevators]string
	buf := make([]byte, 128)
	var State int = 0

	ListenPort, err := net.ResolveUDPAddr("udp", ":20021") // initialize connections
	if err != nil{
		fmt.Println("Error: ", err)
		State = 0
		//time.Sleep(time.Second * 1)
		fmt.Println("First error")
	}
	Listen, err := net.ListenUDP("udp", ListenPort)
	if err != nil{
		fmt.Println("Error: ", err)
		State = 0
		//time.Sleep(time.Second * 1)
		fmt.Println("Second error")
	} else {
		State = 1
		//fmt.Println("State = 1")
	}

	for{
		if State == 0{
			fmt.Println("Looking for connections")
			buffermessage := <- message // message loopback
			recievedmessage <- buffermessage


			RetryListenPort, err := net.ResolveUDPAddr("udp", ":20021")
			if err != nil{
				fmt.Println("Error: ", err)
				State = 0
				time.Sleep(time.Second * 1)
			}
			RetryListen, err := net.ListenUDP("udp", RetryListenPort)
			if err != nil{
				fmt.Println("Error: ", err)
				State = 0
				time.Sleep(time.Second * 1)

			} else {
				State = 1
				fmt.Println("Listener address resolved")
				Listen = RetryListen
			}

		} else if State == 1{
			//fmt.Println("Network working")

			n, addr, err := Listen.ReadFromUDP(buf)

			if err != nil {
				fmt.Println("Error: ", err)
				//Listen.Close()
				//net.ListenUDP("udp", ListenPort)
				State = 0
				fmt.Println("Return to looking for connections")
				//fmt.Println("Connection closed")

			} else {
				//fmt.Println("Recieving messages")
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
							if err != nil{
								fmt.Println("Disconnect :", addressArray[i])
								for j := 0; j <= (numberofIPs - i); j++ {
									addressArray[i+j] = addressArray[i+j+1]
								}
								numberofIPs--
							} else {
								Client, err := net.DialTCP("tcp", nil, Clientaddress)
								if err != nil {
									fmt.Println("Disconnect :", addressArray[i])
									for j := 0; j <= (numberofIPs - i); j++ {
										addressArray[i+j] = addressArray[i+j+1]
									}

									numberofIPs--

								} else {
									_, err = Client.Write([]byte(sendmessage))
									if err != nil{
										fmt.Println("Error: ", err)
									}
								}
							}
						}
					}
				default:
				}
			}
		}
	}
}

func TCP_listener(recievedmessage chan string) {
	listenPort, err := net.Listen("tcp", ":20021")
	if err != nil{
		fmt.Println("Error: ", err)
	}

	for {
		connection, err := listenPort.Accept()
		if err != nil{
			fmt.Println("Error: ", err)
			time.Sleep(time.Second * 1)
			RetrylistenPort, err := net.Listen("tcp", ":20021")
			if err != nil{
				fmt.Println("Error: ", err)
			} else {
				listenPort = RetrylistenPort
			}
		}
		addr := connection.RemoteAddr()

		firsthalf := strings.Split(addr.String(), ":")
		IP := firsthalf[0]

		if IP != "127.0.0.1" {
			buf := make([]byte, 1024)
			n, err := connection.Read(buf)
			if err != nil{
				fmt.Println("Error: ", err)
			}

			address := connection.RemoteAddr().String()
			recievedmessage <- strings.Join([]string{string(buf[0:n]), address}, ",")
			//fmt.Println("Recieved message : ", string(buf[0:n]), " from ", address)
			connection.Close()
		} else {
			connection.Close()
		}
	}

}

