package main

import (
	. "./Heismodul"
	. "./Heismodul/driver"
	. "./Nettverksmodul"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

const N_FLOORS int = 4

func Beregn_kostnad() {
	//gjør beregninger - eksterne og interne ordre
	//det enkleste er å se på det totale antall ordre og etasjer den skal kjøre
	//del resultat på nettverk
}

func Ekstern_ordre(message chan string) {
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
	var currentFloor int
	var floor int
	var newOrder int
	var orderMatch bool
	numberofOrders := [2]int{0, 0}
	var buttonPress [2][4]int
	buttonRelease := [2][4]int{{0, 0, 0, 0}, {0, 0, 0, 0}}
	orderArray := [2][4]int{{-1, -1, -1, -1}, {-1, -1, -1, -1}}
	Button := BUTTON_CALL_UP
	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			currentFloor = floor
		}

		for DIR := 0; DIR <= 1; DIR++ {
			if DIR == 0 {
				Button = BUTTON_CALL_UP
			} else {
				Button = BUTTON_CALL_DOWN
			}

			for FLOOR := 0; FLOOR < N_FLOORS; FLOOR++ { // read buttonpress and add order to que
				buttonPress[DIR][FLOOR] = Elev_get_button_signal(Button, FLOOR) //UP == 0, DOWN == 1
				if (buttonPress[DIR][FLOOR] == 1) && (buttonRelease[DIR][FLOOR] == 0) {
					buttonRelease[DIR][FLOOR] = 1
					fmt.Println("New buttonpress at: ", FLOOR)
					orderMatch = false
					newOrder = FLOOR
					for j := 0; j <= numberofOrders[DIR]; j++ {
						if orderArray[DIR][j] == newOrder {
							orderMatch = true
							fmt.Println("Order already exist")
						}
					}

					if (orderMatch == false) && (currentFloor != newOrder) {
						orderArray[DIR][numberofOrders[DIR]] = newOrder
						fmt.Println("New order at floor: ", newOrder)
						//floor := strconv.FormatInt(int64(newOrder), 10)
						//direction := strconv.FormatInt(int64(DIR), 10)
						//call := []string{floor, direction}
						message <- strings.Join([]string{strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
						Elev_set_button_lamp(Button, orderArray[DIR][numberofOrders[DIR]], 1) // Flyttes etterhvert
						numberofOrders[DIR]++
					}
				} else if (buttonPress[DIR][FLOOR] == 0) && (buttonRelease[DIR][FLOOR] == 1) {
					//fmt.Println("New buttonrelease at: ", i)
					buttonRelease[DIR][FLOOR] = 0
				}
			}
		}
	}
}

func Intern_ordre(nextFloor chan int, orderFinished chan bool) {

	var currentFloor = Init_floor()

	go Displayfloor()
	go Kjør_heis(nextFloor, orderFinished)

	var floor int
	var numberofOrders int = 0
	var orderArray [N_FLOORS + 1]int //initialize orderArray
	for j := 0; j <= N_FLOORS; j++ {
		orderArray[j] = -1
	}

	var newOrder int
	var orderMatch bool
	var buttonPress [4]int
	buttonRelease := [4]int{0, 0, 0, 0}

	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			currentFloor = floor
		}

		if numberofOrders > 0 { // go to floor and remove order from que when finished

			select {
			case orderFinished_i := <-orderFinished:
				if orderFinished_i == true {
					Elev_set_button_lamp(BUTTON_COMMAND, orderArray[0], 0)
					fmt.Println("Order to floor: ", orderArray[0], " finished, removed from que")
					for i := 0; i < numberofOrders; i++ {
						orderArray[i] = orderArray[i+1]
					}
					numberofOrders--
					fmt.Println("Number of orders: ", numberofOrders)
					if numberofOrders >= 1 {
						nextFloor <- orderArray[0]
						fmt.Println("Next floor is: ", orderArray[0])
					}
				}
			default:
			}
		}

		for i := 0; i < N_FLOORS; i++ { // read buttonpress and add order to que
			buttonPress[i] = Elev_get_button_signal(BUTTON_COMMAND, i)
			if (buttonPress[i] == 1) && (buttonRelease[i] == 0) {
				buttonRelease[i] = 1
				fmt.Println("New buttonpress at: ", i)
				orderMatch = false
				newOrder = i
				for j := 0; j <= numberofOrders; j++ {
					if orderArray[j] == newOrder {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if (orderMatch == false) && (currentFloor != newOrder) {
					orderArray[numberofOrders] = newOrder
					fmt.Println("New order at floor: ", newOrder)
					Elev_set_button_lamp(BUTTON_COMMAND, orderArray[numberofOrders], 1)
					numberofOrders++
					if numberofOrders == 1 {
						nextFloor <- orderArray[0]
						fmt.Println("Next floor is: ", orderArray[0])
					}
				}
			} else if (buttonPress[i] == 0) && (buttonRelease[i] == 1) {
				//fmt.Println("New buttonrelease at: ", i)
				buttonRelease[i] = 0
			}
		}
	}
}

func Inkommende_ordre(recievedmessage chan string, nextFloor chan int) {
	//check om ordre allerede er i køen
	//legg til hvis ikke
	//fjern intern ordre fra kø når ordren er fullført
	for {
		select {
		case newOrder := <-recievedmessage:
			{
				//fmt.Print(newOrder)
				slice := strings.Split(newOrder, ",")
				//var first int = int(slice[0])
				//fmt.Println(first)
				value, err := strconv.ParseInt(slice[0], 10, 64)
				CheckError(err)
				nextFloor <- int(value)
			}
		default:
		}
	}
}

func Vurder_kostnad() {
	//sammenlign innkommende resultat
	//Vurder om vi skal ta ordre og legge den til i intern ordrekø
}

func Kvitter_ordre() {
	//send kvittering for ekstern ordre på nettverket
	//motta kvittering
	//fjern ordre fra kø
}

func main() {
	nextFloor := make(chan int, 10)
	orderFinished := make(chan bool, 1)
	message := make(chan string, 1024)
	recievedmessage := make(chan string, 1024)
	runtime.GOMAXPROCS(runtime.NumCPU())

	go Intern_ordre(nextFloor, orderFinished)
	go Ekstern_ordre(message)
	go Broadcast()
	go TCP_sender(message)
	go TCP_listener(recievedmessage)
	go Inkommende_ordre(recievedmessage, nextFloor)
	deadChan := make(chan bool, 1)
	<-deadChan
}
