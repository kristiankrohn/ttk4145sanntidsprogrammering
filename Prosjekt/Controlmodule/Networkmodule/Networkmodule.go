package Networkmodule

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//messages on the network is sent as strings

const Numberofelevators int = 10

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func Last_byte_of_my_IP() int { //was borrowed from https://github.com/TTK4145/Network-go/blob/master/network/localip/localip.go, but rewritten to the unreconizeable
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
	lastaddressbyte_i64, err := strconv.ParseInt(lastaddressbytestring[3], 10, 64)
	CheckError(err)
	lastaddressbyte := int(lastaddressbyte_i64)
	return lastaddressbyte
}

func Broadcast(message chan string, recievedmessage chan string) {
	//Expose ourselves to other elevators, if network is not working loopback messages to ourselves
	for {
		BroadcastAddr, err := net.ResolveUDPAddr("udp", "129.241.187.255:20021")
		if err != nil {
			fmt.Println("Warning: ", err)
		}
		BroadcastSpammer, err := net.DialUDP("udp", nil, BroadcastAddr)
		if err != nil {
			fmt.Println("Warning: ", err)
			select {
			case buffermessage := <-message:
				{ // message loopback, remove when TCP_sender throws an no connection erre
					recievedmessage <- strings.Join([]string{buffermessage, "0.0.0.255:9000"}, ",")
				}
			default:
			}

		} else {
			BroadcastSpammer.Write([]byte("Hello, i'm an elevator - connect to me!"))
			BroadcastSpammer.Close()
		}
		time.Sleep(time.Second * 1)
	}
}

func TCP_sender(message chan string, recievedmessage chan string) {
	//Discover other elevators and put in elevatorarray, when message is recieved - send to all elevators
	var numberofIPs int = 0
	var addressArray [Numberofelevators]string
	buf := make([]byte, 128)
	var State int = 0

	ListenPort, err := net.ResolveUDPAddr("udp", ":20021") // initialize connections
	if err != nil {
		fmt.Println("Error: ", err)
		State = 0
		fmt.Println("First TCP error")
	}
	Listen, err := net.ListenUDP("udp", ListenPort)
	if err != nil {
		fmt.Println("Error: ", err)
		State = 0
		fmt.Println("Second TCP error")
	} else {
		State = 1
	}

	for {
		if State == 0 {
			fmt.Println("Looking for connections")
			select {
			case buffermessage := <-message:
				{ // message loopback, program has never made it here during testing, but preferably message loopback would have happened here(not that it functionally matters)
					recievedmessage <- strings.Join([]string{buffermessage, "0.0.0.255:9000"}, ",")
				}
			default:
			}

			RetryListenPort, err := net.ResolveUDPAddr("udp", ":20021")
			if err != nil {
				fmt.Println("Error: ", err)
				State = 0
				time.Sleep(time.Second * 1)
			}
			RetryListen, err := net.ListenUDP("udp", RetryListenPort)
			if err != nil {
				fmt.Println("Error: ", err)
				State = 0
				time.Sleep(time.Second * 1)

			} else {
				State = 1
				fmt.Println("Listener address resolved")
				Listen = RetryListen
			}

		} else if State == 1 {
			//Discover new elevators

			n, addr, err := Listen.ReadFromUDP(buf)

			if err != nil {
				fmt.Println("Error: ", err)
				State = 0
				fmt.Println("Return to looking for connections")

			} else {
				firsthalf := strings.Split(addr.String(), ":")
				NewIP := firsthalf[0]
				var IPmatch bool = false

				if string(buf[0:n]) == "Hello, i'm an elevator - connect to me!" {

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

				select { //message sending happens here
				case sendmessage := <-message:
					{
						for i := 0; i <= numberofIPs; i++ {
							Clientaddress, err := net.ResolveTCPAddr("tcp", string(addressArray[i]+":20021"))
							if err != nil {
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
									if err != nil {
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
		time.Sleep(time.Millisecond * 10)

	}
}

func TCP_listener(recievedmessage chan string) {
	//Recieves messages and put on channel
	listenPort, err := net.Listen("tcp", ":20021")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for {
		connection, err := listenPort.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			time.Sleep(time.Second * 1)
			RetrylistenPort, err := net.Listen("tcp", ":20021")
			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				listenPort = RetrylistenPort
			}
		}
		addr := connection.RemoteAddr()

		firsthalf := strings.Split(addr.String(), ":")
		IP := firsthalf[0]

		//We do not want messages from 127.0.0.1!
		if IP != "127.0.0.1" {
			buf := make([]byte, 1024)
			n, err := connection.Read(buf)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			address := connection.RemoteAddr().String()
			recievedmessage <- strings.Join([]string{string(buf[0:n]), address}, ",")
			connection.Close()
		} else {
			connection.Close()
		}
		time.Sleep(time.Millisecond * 10)

	}
}
