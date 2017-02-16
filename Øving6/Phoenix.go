package main

import (
  "net"
  "fmt"
  "time"
  "os/exec"
  "os"
  "encoding/binary"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func risefromtheashes()(uint64){

  var counter uint64 = 0
  buffer := make([]byte, 8)
  ListenAddr, _ := net.ResolveUDPAddr("udp", "129.241.187.255:20021")
	Listener, err := net.ListenUDP("udp", ListenAddr)
  checkError(err)

  fmt.Println("Connection established")
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

  command := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run Phoenix.go")
  _ = command.Run()
  fmt.Println("The Phoenix has risen")

  return counter
}

func main(){
  var counter uint64 = risefromtheashes()
  isAliveAddr, err := net.ResolveUDPAddr("udp", "129.241.187.255:20021")
	checkError(err)
	isAlive, err := net.DialUDP("udp", nil, isAliveAddr)
	checkError(err)
  for{
    counter++
    buffer := make([]byte, 8)
    binary.BigEndian.PutUint64(buffer, counter)
    _,_ = isAlive.Write(buffer)
  //  checkError(err)
    fmt.Println(counter)
    time.Sleep(time.Second * 1)
  }
}
