package main

import (
	. "./Heismodul"
	. "./Heismodul/driver"
	. "./Nettverksmodul"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//const N_FLOORS int = 4 Define this in Heismodul.go
//const numberofelevators int = 5 //must be higher than maximum number of possible elevators, or it will cause bufferoverflow
const arraysize int = N_FLOORS * 3 // number of buttons as a function of number of elevators

var orderArray [2][arraysize]int
var ext_orderArray [2][arraysize]int
var cost_array[arraysize][Numberofelevators] Costentry
var numberofCosts[arraysize] Costnumber
var numberofOrders int
var ext_numberofOrders int
//var currentFloor int - Declared in Heismodul

type Costentry struct{
	cost int
	IP int
}

type Costnumber struct{
	number int
	starttime time.Time
}


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

	for i := 0; i < 2; i++ {
		for j := 0; j < arraysize; j++ {
			ext_orderArray[i][j] = -1
		}
	}

	for i := 0; i < arraysize; i++{
		for j := 0; j < Numberofelevators; j++{
			cost_array[i][j] = Costentry{0, 0}
		}
	}
	initCostnumber := Costnumber{0, time.Now()}
	for k := 0; k < arraysize; k++{
		numberofCosts[k] = initCostnumber
	}
	numberofOrders = 0
	//Init_elevator()

}

func Calculate_cost(floor int, calldirection int) int{
	//gjør beregninger - eksterne og interne ordre
	//det enkleste er å se på det totale antall ordre og etasjer den skal kjøre
	//del resultat på nettverk
	var cost int
	var direction int = 1000

	if floor > CurrentFloor{
		if orderArray[0][0] > CurrentFloor{
			direction = 1
		} else {
			direction = 2
		}
	} else {
		if orderArray[0][0] > CurrentFloor{
			direction = 2
		} else {
			direction = 1
		}
	}

	cost = direction * (floor - CurrentFloor)

	return cost
}

func Local_orders(internal_button chan int, nextFloor chan int, orderFinished chan bool) {


	//go Displayfloor()
	//go Elevator_driver(nextFloor, orderFinished)
	//go Handle_buttons(up_button, down_button, internal_button)
	var orderMatch bool = false

	for {
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

func External_orders(message chan string, up_button chan int, down_button chan int, nextFloor chan int) {
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

	var newOrder int
	var orderMatch bool
	var DIR = 0

	for {

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
					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
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
					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
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

func Incomming_message(recievedmessage chan string, message chan string) {
	//check om ordre allerede er i køen
	//legg til hvis ikke
	//fjern intern ordre fra kø når ordren er fullført
	
	for {
		select {
		case newOrder := <-recievedmessage:
			{	
				//Messagecodes: 
				// 0 = new order
				// 1 = cost
				// 3 = kvittering
				//fmt.Print(newOrder)
				slice := strings.Split(newOrder, ",")
				//var first int = int(slice[0])
				//fmt.Println(first)
				messagecode , err := strconv.ParseInt(slice[0], 10, 64)
				CheckError(err)

				if messagecode == 0{
					floor, err := strconv.ParseInt(slice[1], 10, 64)
					CheckError(err)
					newOrder := int(floor)
					direction, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)
				//nextFloor <- int(value)
					DIR := int(direction)
					if DIR == 0 {
						ext_orderArray[0][numberofOrders] = newOrder
						ext_orderArray[1][numberofOrders] = DIR
						Elev_set_button_lamp(BUTTON_CALL_UP, ext_orderArray[0][numberofOrders], 1) // Flyttes etterhvert
						ext_numberofOrders++

						cost := Calculate_cost(newOrder, DIR)
						button := newOrder * (DIR + 1)

						message <- strings.Join([]string{strconv.FormatInt(int64(1), 10), strconv.FormatInt(int64(cost), 10), strconv.FormatInt(int64(button), 10)}, ",") 
					/*if numberofOrders == 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
						}*/
					} else if DIR == 1 {
						ext_orderArray[0][numberofOrders] = newOrder
						ext_orderArray[1][numberofOrders] = DIR
						Elev_set_button_lamp(BUTTON_CALL_DOWN, orderArray[0][numberofOrders], 1) // Flyttes etterhvert
						ext_numberofOrders++

						cost := Calculate_cost(newOrder, DIR)
						
						button := newOrder * (DIR + 1)
						message <- strings.Join([]string{strconv.FormatInt(int64(1), 10), strconv.FormatInt(int64(cost), 10), strconv.FormatInt(int64(button), 10)}, ",") 
					/*if numberofOrders == 1 {
						nextFloor <- orderArray[0][0]
						fmt.Println("Next floor is: ", orderArray[0][0])
						}*/
					}
				} else if messagecode == 1{
					cost_s, err := strconv.ParseInt(slice[1], 10, 64)
					CheckError(err)
					cost := int(cost_s)
					button_s, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)
					button := int(button_s)
					lastaddressbytestring_pp := strings.Split(slice[3], ":")
					lastaddressbytestring := strings.Split(lastaddressbytestring_pp[0], ".")
					lastaddressbyte_i64, err:= strconv.ParseInt(lastaddressbytestring[3], 10, 64)
					CheckError(err)
					lastaddressbyte := int(lastaddressbyte_i64)
					fmt.Println(lastaddressbyte)
					
					newcost := Costentry{cost, lastaddressbyte}

					cost_array[button][numberofCosts[button].cost] = newcost
					if numberofCosts[button].number == 0{
						numberofCosts[button].starttime = time.Now()
					}
					numberofCosts[button].number ++
					
				} else if messagecode == 2{

				}
			}
		default:
		}
	}
}

func Assess_cost() {
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
	//current_floor_internal := make(chan int, 10)
	//current_floor_external := make(chan int, 10)
	message := make(chan string, 1024)
	recievedmessage := make(chan string, 1024)
	runtime.GOMAXPROCS(runtime.NumCPU())
	Init_system()
	fmt.Println("Init finished")
	go Elevator_driver(nextFloor, orderFinished)

	go Local_orders(internal_button, nextFloor, orderFinished)
	go Handle_buttons(up_button, down_button, internal_button)
	go External_orders(message, up_button, down_button, nextFloor)
	//go Current_floor()
	go Displayfloor()
	go Broadcast(message, recievedmessage)
	go TCP_sender(message, recievedmessage)
	go TCP_listener(recievedmessage)
	go Incomming_message(recievedmessage, message)
	go Clear_orders(orderFinished, nextFloor)

	deadChan := make(chan bool, 1)
	<-deadChan
}
