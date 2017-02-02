package main

import (
    "fmt"
    "net"
    "os"
		"time"
    "strings"
    //"bytes"
)
const numberofelevators int = 4
//var global_addressArray [numberofelevators]string
//var global_numberofIPs int = 0
/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

func Broadcast(){

  BroadcastAddr, err := net.ResolveUDPAddr("udp", "129.241.187.255:20021")
  CheckError(err)
  BroadcastSpammer, err := net.DialUDP("udp", nil, BroadcastAddr)
  CheckError(err)

  for{
    BroadcastSpammer.Write([]byte("Hello, i'm an elevator from group 67 - connect to me!"))
    time.Sleep(time.Second * 1)
  }
}

func Listener(){
  ListenPort, err := net.ResolveUDPAddr("udp",":20021")
  CheckError(err)
  Listen, err := net.ListenUDP("udp", ListenPort)
  CheckError(err)
  buf := make([]byte, 1024)

  var numberofIPs int = 0
  var addressArray [numberofelevators]string
  var IPmatch bool
	for {
			n,addr,err := Listen.ReadFromUDP(buf)
			//fmt.Println("Received ",string(buf[0:n]), " from ",addr)

			if err != nil {
					fmt.Println("Error: ",err)
			}
      //IPstring := addr.String()
      firsthalf := strings.Split(addr.String(), ":")
      NewIP := firsthalf[0]

      if (string(buf[0:n]) == "Hello, i'm an elevator from group 67 - connect to me!"){
        IPmatch = false
        //numberofIPs <- global_numberofIPs
        for i := 0; i <= numberofIPs; i++ {

            if addressArray[i] == NewIP{
                IPmatch = true

            }
          }

          if IPmatch == false{
              //numberofIPs <- global_numberofIPs
              //addressArray <- global_addressArray

              numberofIPs ++
              addressArray[numberofIPs] = NewIP

              //global_numberofIPs <- numberofIPs
              //global_addressArray <- addressArray

              fmt.Println("New machine at: ", addressArray[numberofIPs])
            }
      }

  }
}

func main()  {


  go Broadcast()
  go Listener()

  deadChan :=make(chan bool, 1)
  <- deadChan
}
