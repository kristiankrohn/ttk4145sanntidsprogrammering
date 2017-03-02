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

//const N_FLOORS int = 4
const arraysize int = 12

var orderArray [2][arraysize]int
var numberofOrders int

/* 	direction up = 0
*	direction down = 1
*	direction internal = 2
 */

func Init_system() {
	Init_elevator()
	for i := 0; i < 2; i++ {
		for j := 0; j < arraysize; j++ {
			orderArray[i][j] = -1
		}
	}
	numberofOrders = 0
	//Init_elevator()

}

func Beregn_kostnad() {
	//gjør beregninger - eksterne og interne ordre
	//det enkleste er å se på det totale antall ordre og etasjer den skal kjøre
	//del resultat på nettverk
}

func Local_orders(internal_button chan int, current_floor_internal chan int, nextFloor chan int, orderFinished chan bool) {

	var currentFloor int
	//go Displayfloor()
	//go Elevator_driver(nextFloor, orderFinished)
	//go Handle_buttons(up_button, down_button, internal_button)
	var orderMatch bool = false

	for {
		select {
		case floor := <-current_floor_internal:
			{
				currentFloor = int(floor)
				//fmt.Println("Current floor is: ", currentFloor)
				currentFloor = currentFloor
			}
		default:
		}

		select {
		case internal_call := <-internal_button:
			{
				newOrder := int(internal_call)
				for j := 0; j <= numberofOrders; j++ {
					if (orderArray[0][j] == newOrder) && (orderArray[1][j] == 2) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {
					//if (orderMatch == false) && (currentFloor != newOrder) {
					//fmt.Println("Current floor is: ", currentFloor)
					orderArray[0][numberofOrders] = newOrder
					orderArray[1][numberofOrders] = 2
					fmt.Println("New order at floor: ", newOrder)
					Elev_set_button_lamp(BUTTON_COMMAND, orderArray[0][numberofOrders], 1)
					numberofOrders++
					if numberofOrders == 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
					}
				}
			}

		default:
		}
	}
}

func External_orders(message chan string, up_button chan int, down_button chan int, current_floor_external chan int, nextFloor chan int) {
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
	var newOrder int
	var orderMatch bool
	var DIR = 0

	for {

		select {
		case floor := <-current_floor_external:
			{
				currentFloor = int(floor)
				currentFloor = currentFloor
			}
		default:
		}

		select {
		case call_up := <-up_button:
			{
				DIR = 0
				newOrder = int(call_up)
				for j := 0; j <= numberofOrders; j++ {
					if (orderArray[0][j] == newOrder) && (orderArray[1][j] == DIR) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {
					//orderArray[0][numberofOrders] = newOrder
					//orderArray[1][numberofOrders] = DIR
					fmt.Println("New order at floor: ", newOrder)
					//floor := strconv.FormatInt(int64(newOrder), 10)
					//direction := strconv.FormatInt(int64(DIR), 10)
					//call := []string{floor, direction}
					message <- strings.Join([]string{strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
					/*Elev_set_button_lamp(BUTTON_CALL_UP, orderArray[0][numberofOrders], 1) // Flyttes etterhvert
					numberofOrders++
					if numberofOrders == 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
					}*/
				}
			}
		default:
		}

		select {
		case call_down := <-down_button:
			{
				DIR = 1
				newOrder = int(call_down)
				for j := 0; j <= numberofOrders; j++ {
					if (orderArray[0][j] == newOrder) && (orderArray[1][j] == DIR) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {
					//orderArray[0][numberofOrders] = newOrder
					//orderArray[1][numberofOrders] = DIR
					fmt.Println("New order at floor: ", newOrder)
					//floor := strconv.FormatInt(int64(newOrder), 10)
					//direction := strconv.FormatInt(int64(DIR), 10)
					//call := []string{floor, direction}
					message <- strings.Join([]string{strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
					/*Elev_set_button_lamp(BUTTON_CALL_DOWN, orderArray[0][numberofOrders], 1) // Flyttes etterhvert
					numberofOrders++
					if numberofOrders == 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
					}*/
				}
			}
		default:
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
				floor, err := strconv.ParseInt(slice[0], 10, 64)
				CheckError(err)
				newOrder := int(floor)
				direction, err := strconv.ParseInt(slice[1], 10, 64)
				CheckError(err)
				//nextFloor <- int(value)
				DIR := int(direction)
				if DIR == 0 {
					orderArray[0][numberofOrders] = newOrder
					orderArray[1][numberofOrders] = DIR
					Elev_set_button_lamp(BUTTON_CALL_UP, orderArray[0][numberofOrders], 1) // Flyttes etterhvert
					numberofOrders++
					if numberofOrders == 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
					}
				} else if DIR == 1 {
					orderArray[0][numberofOrders] = newOrder
					orderArray[1][numberofOrders] = DIR
					Elev_set_button_lamp(BUTTON_CALL_DOWN, orderArray[0][numberofOrders], 1) // Flyttes etterhvert
					numberofOrders++
					if numberofOrders == 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
					}
				}
			}
		default:
		}
	}
}

func Vurder_kostnad() {
	//sammenlign innkommende resultat
	//Vurder om vi skal ta ordre og legge den til i intern ordrekø
}

func Clear_orders(orderFinished chan bool, nextFloor chan int) {
	//send kvittering for ekstern ordre på nettverket
	//motta kvittering
	//fjern ordre fra kø

	Button := BUTTON_COMMAND
	for {
		if numberofOrders > 0 { // go to floor and remove order from que when finished

			select {
			case orderFinished_i := <-orderFinished:
				if orderFinished_i == true {
					if orderArray[1][0] == 0 {
						Button = BUTTON_CALL_UP
					} else if orderArray[1][0] == 1 {
						Button = BUTTON_CALL_DOWN
					} else {
						Button = BUTTON_COMMAND
					}
					Elev_set_button_lamp(Button, orderArray[0][0], 0)
					fmt.Println("Order to floor: ", orderArray[0][0], " finished, removed from que")
					for i := 0; i < numberofOrders; i++ {
						orderArray[0][i] = orderArray[0][i+1]
						orderArray[1][i] = orderArray[1][i+1]
					}
					numberofOrders--
					fmt.Println("Number of orders: ", numberofOrders)
					if numberofOrders >= 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
					}
				}
			default:
			}
		}
	}
}

func main() {

	nextFloor := make(chan int, 10)
	orderFinished := make(chan bool, 1)
	up_button := make(chan int, 10)
	down_button := make(chan int, 10)
	internal_button := make(chan int, 10)
	current_floor_internal := make(chan int, 10)
	current_floor_external := make(chan int, 10)
	message := make(chan string, 1024)
	recievedmessage := make(chan string, 1024)
	runtime.GOMAXPROCS(runtime.NumCPU())
	Init_system()
	fmt.Println("Init finished")
	go Elevator_driver(nextFloor, orderFinished)

	go Local_orders(internal_button, current_floor_internal, nextFloor, orderFinished)
	go Handle_buttons(up_button, down_button, internal_button)
	go External_orders(message, up_button, down_button, current_floor_external, nextFloor)
	go Current_floor(current_floor_external, current_floor_internal)
	go Displayfloor(current_floor_external, current_floor_internal)
	go Broadcast()
	go TCP_sender(message)
	go TCP_listener(recievedmessage)
	go Inkommende_ordre(recievedmessage, nextFloor)
	go Clear_orders(orderFinished, nextFloor)
	deadChan := make(chan bool, 1)
	<-deadChan
}
