package main

import (
    "fmt"
    "net"
    "os"
		"time"
)
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func Read(connection net.Conn){
  buf := make([]byte, 1024)

	for {
			n,err := connection.Read(buf)
			fmt.Println(string(buf[0:n]))

			if err != nil {
					fmt.Println("Error: ",err)
			}
	}
}

func Write(connection net.Conn)  {
  for {
    connection.Write([]byte("Hello from client21!\x00"))
    time.Sleep(time.Second * 1)
  }
}

func main(){
  ServerAddr,err := net.ResolveTCPAddr("tcp","129.241.187.43:33546")
  CheckError(err)

  ClientRead,err := net.DialTCP("tcp", nil, ServerAddr)
  CheckError(err)

  go Read(ClientRead)
  go Write(ClientRead)

  deadChan :=make(chan bool, 1)
  <- deadChan
}
