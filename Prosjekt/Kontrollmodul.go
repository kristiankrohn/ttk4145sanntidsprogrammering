package main

import (
  "fmt"

)

func Beregn_kostnad(){
  //gjør beregninger - eksterne og interne ordre
    //det enkleste er å se på det totale antall ordre og etasjer den skal kjøre
  //del resultat på nettverk
}

func Ekstern_ordre(){
  //read buttonpress
  //check om ordre allerede er i køen
  //legg til hvis ikke
  //del ordre på nettverket
  //del på nytt hvis timeout

/************************************************
Sålangt:
Prøver å lese buttonpress up & down,
på samme måte som med internal buttons.
To av alt - opp og ned
Tar i mot ordre og lagrer i array, men
gjør ingenting annet enda. Sjekker om
ordre allerede finnes.

Trenger å dele ordre på nettverket og vente
på svar. Eventuelt sende på nytt hvis timeout.
Etter at noen har tatt ordren kan lampen tennes.
************************************************/

  var numberofOrdersUp int = 0
  var numberofOrdersUp int = 0

  var orderArrayUp [N_FLOORS + 1]int
  var orderArrayDown [N_FLOORS + 1]int
	for j := 0; j <= N_FLOORS; j++ {
		orderArrayUp[j] = -1
    orderArrayDown[j] = -1
	}

  var newExternalOrder int
  var orderMatchUp bool
  var buttonPressUp [3]int
  buttonReleaseUp := [3]int{0, 0, 0}
  buttonReleaseDown := [3]int{0, 0, 0}

  for i := 0; i < N_FLOORS; i++ {
    buttonPressUp[i] = Elev_get_button_signal(BUTTON_CALL_UP, i)
    buttonPressDown[i] = Elev_get_button_signal(BUTTON_CALL_DOWN, i)

    if (buttonPressUp [i] == 1) && (buttonReleaseUp[i] == 0) {
      buttonReleaseUp[i] = 1
      fmt.Println("New buttonpress up at: ", i)
      orderMatchUp = false
      newOrderUp = i
      for j := 0; j <= numberofOrdersUp; j++ {
        if orderArrayUp[j] == newOrderUp {
          orderMatchUp = true
          fmt.Println("Order already exist")
          //ÅPNE DØR??
        }
      }
      if (orderMatchUp == false) {
        //SEND PÅ NETTVERK
        orderArrayUp[numberofOrdersUp] = newOrderUp
        fmt.Println("New order up at floor: ", newOrderUp)
        Elev_set_button_lamp(BUTTON_CALL_UP, orderArrayUP[numberofOrdersUp], 1)
        numberofOrdersUp++
      }

    } else if (buttonPressUp[i] == 0) && (buttonReleaseUp[i] == 1) {
      //fmt.Println("New buttonrelease at: ", i)
      buttonReleaseUp[i] = 0
    }

    if (buttonPressDown [i] == 1) && (buttonReleaseDown[i] == 0) {
      buttonReleaseDown[i] = 1
      fmt.Println("New buttonpress down at: ", i)
      orderMatchDown = false
      newOrderDown = i
      for j := 0; j <= numberofOrdersDownn; j++ {
        if orderArrayDown[j] == newOrderDown {
          orderMatchDown = true
          fmt.Println("Order already exist")
          //ÅPNE DØR??
        }
      }
      if (orderMatchDown == false) {
        //SEND PÅ NETTVERK
        orderArrayDown[numberofOrdersDown] = newOrderDown
        fmt.Println("New order down at floor: ", newOrderDown)
        Elev_set_button_lamp(BUTTON_CALL_Down, orderArrayDown[numberofOrdersDown], 1)
        numberofOrdersDown++
      }

    } else if (buttonPressDown[i] == 0) && (buttonReleaseDown[i] == 1) {
      //fmt.Println("New buttonrelease at: ", i)
      buttonReleaseDown[i] = 0
    }
  }
}

func Inkommende_ordre(){
  //check om ordre allerede er i køen
  //legg til hvis ikke
  //fjern intern ordre fra kø når ordren er fullført
}

func Vurder_kostnad(){
  //sammenlign innkommende resultat
  //Vurder om vi skal ta ordre og legge den til i intern ordrekø
}

func Kvitter_ordre(){
  //send kvittering for ekstern ordre på nettverket
  //motta kvittering
  //fjern ordre fra kø
}

func main(){


  deadChan := make(chan bool, 1)
	<-deadChan
}
