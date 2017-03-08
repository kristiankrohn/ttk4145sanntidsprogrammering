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
const arraysize int = N_FLOORS * 10 // number of buttons as a function of number of elevators

var orderArray [2][arraysize]int
var ext_orderArray [2][arraysize] Extentry
var cost_array[arraysize][Numberofelevators] Costentry
var numberofCosts[arraysize] Costnumber
var numberofOrders int
var ext_numberofOrders int
//var currentFloor int - Declared in Heismodul.go
type Orderentry struct{
	floor int
	button int
}

type Costentry struct{
	cost int
	IP int
}

type Costnumber struct{
	number int
	starttime time.Time
}

type Extentry struct{
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
			ext_orderArray[i][j] = Extentry{-1, time.Now()}
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
	var direction int = 10

	if numberofOrders == 0{
		if floor > CurrentFloor{
			cost = floor - CurrentFloor
		} else {
			cost = CurrentFloor - floor
		}
	} else { 
		if floor > CurrentFloor{
			if orderArray[0][0] > CurrentFloor{
				direction = 1
			} else {
				direction = 2
			}
			cost = direction * (floor - CurrentFloor)
		} else {
			if orderArray[0][0] > CurrentFloor{
				direction = 2
			} else {
				direction = 1
			}
			cost = direction * (CurrentFloor - floor)
		}
	}
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
				orderMatch = false
				DIR = 0
				newOrder = int(call_up)
				for j := 0; j <= ext_numberofOrders; j++ {
					if (ext_orderArray[0][j].number == newOrder) && (ext_orderArray[1][j].number == DIR) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {

					fmt.Println("New order at floor: ", newOrder)

					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
				}
			}
		default:
		}

		select {
		case call_down := <-down_button:
			{
				orderMatch = false
				DIR = 1
				newOrder = int(call_down)
				for j := 0; j <= ext_numberofOrders; j++ {
					if (ext_orderArray[0][j].number == newOrder) && (ext_orderArray[1][j].number == DIR) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {

					fmt.Println("New order at floor: ", newOrder)

					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
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

				messagecode , err := strconv.ParseInt(slice[0], 10, 64)
				CheckError(err)

				if messagecode == 0{ // Recieving a new order and sening out the cost
					floor, err := strconv.ParseInt(slice[1], 10, 64)
					CheckError(err)
					newOrder := int(floor)
					direction, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)

					DIR := int(direction)
					if ext_numberofOrders < 0 {ext_numberofOrders = 0}

					if DIR == 0 {

						Elev_set_button_lamp(BUTTON_CALL_UP, newOrder, 1) // Flyttes etterhvert?				
					} else if DIR == 1 {
					
						Elev_set_button_lamp(BUTTON_CALL_DOWN, newOrder, 1) // Flyttes etterhvert?	
					}

					orderMatch := false
					
					for j := 0; j <= ext_numberofOrders; j++ {
						if (ext_orderArray[0][j].number == newOrder) && (ext_orderArray[1][j].number == DIR) {
							orderMatch = true
							fmt.Println("Extorder already exist")
						}
					}
					if orderMatch == false{ // we have a new order that doesn't exist in the array
						ext_orderArray[0][ext_numberofOrders] = Extentry{newOrder, time.Now()}
						ext_orderArray[1][ext_numberofOrders] = Extentry{DIR, time.Now()}
						ext_numberofOrders++

						fmt.Println("Remaining external order: ", ext_numberofOrders)
						fmt.Println("New external order array", ext_orderArray[0][0].number)
					}

					cost := Calculate_cost(newOrder, DIR)

					button := newOrder + (DIR+1) * (DIR+1)
					
					message <- strings.Join([]string{strconv.FormatInt(int64(1), 10), strconv.FormatInt(int64(cost), 10), strconv.FormatInt(int64(button), 10)}, ",") 
					fmt.Println("Returning Cost: ", cost)
					
						
				} else if messagecode == 1{ // Gather costs and put them in costarray
					// Slicing string and convert to useful datatypes
					cost_i64, err := strconv.ParseInt(slice[1], 10, 64) // converting via i64
					CheckError(err)
					cost := int(cost_i64)
					button_i64, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)
					button := int(button_i64)
					lastaddressbytestring_pp := strings.Split(slice[3], ":")
					lastaddressbytestring := strings.Split(lastaddressbytestring_pp[0], ".")
					lastaddressbyte_i64, err:= strconv.ParseInt(lastaddressbytestring[3], 10, 64)
					CheckError(err)
					lastaddressbyte := int(lastaddressbyte_i64)
					fmt.Println("Adding cost to array", cost, button, numberofCosts[button].number)//Cost_array = arraysize * numberofelevators					
					
					// Adding cost to costarray and starting timer if it is first cost added
					cost_array[button][numberofCosts[button].number] = Costentry{cost, lastaddressbyte}
					if numberofCosts[button].number == 0{
						numberofCosts[button].starttime = time.Now()
					}
					numberofCosts[button].number ++
					
				} else if messagecode == 2{ //Ext_order is completed and removed from ext_array
					floor_i64, err := strconv.ParseInt(slice[1], 10, 64)
					CheckError(err)

					DIR_i64, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)

					floor := int(floor_i64)
					DIR := int(DIR_i64)

					Button := BUTTON_CALL_UP
					if DIR == 0{
						Button = BUTTON_CALL_UP
					} else {
						Button = BUTTON_CALL_DOWN
					}

					Elev_set_button_lamp(Button, floor, 0)

					for i := 0; i < ext_numberofOrders; i++{
						if (ext_orderArray[0][i].number == floor) && (ext_orderArray[1][i].number == DIR){
							for j:= i; j < ext_numberofOrders; j++{
								ext_orderArray[0][j] = ext_orderArray[0][j + 1]
								ext_orderArray[1][j] = ext_orderArray[1][j + 1]
								i = j
							}
							ext_numberofOrders --
							fmt.Print("Removed completed externalorder, remaining: ", ext_numberofOrders)
							fmt.Println(" First is to floor: ", ext_orderArray[0][0].number)
						}
					}

					//ext_numberofOrders --
					if ext_numberofOrders < 0 {ext_numberofOrders = 0} // just to be sure 
				}
			}
		default:
		}
	}
}

func Assess_cost(nextFloor chan int) {
	//sammenlign innkommende resultat
	//Vurder om vi skal ta ordre og legge den til i intern ordrekø
	for{
		myIP := Last_byte_of_my_IP()
		for i := 0; i < arraysize; i++{
			now := time.Now()
			//fmt.Println("checking for ordertimeouts: ", i)
			if (now.Sub(numberofCosts[i].starttime) > 1000000000) && (numberofCosts[i].number > 0){ // check for timeout, if timeout, assess costarray
				var min Costentry = cost_array[i][0]
				fmt.Println("Order auction ended, assessing cost")
				for j := 0; j < numberofCosts[i].number; j++{
					if min.cost > cost_array[i][j].cost{
						min.cost = cost_array[i][j].cost
					}

				}
				numberofCosts[i].number = 0
				if (myIP <= min.IP) {
					fmt.Println("I have lowest cost and taking order, the number in que is", numberofOrders)
					if i < 4{
						orderArray[0][numberofOrders] = i - 1
						orderArray[1][numberofOrders] = 0

					} else {
						orderArray[0][numberofOrders] = i - 4
						orderArray[1][numberofOrders] = 1

					}
					numberofOrders ++

					if numberofOrders == 1 {
						fmt.Println("This is first order and sending to elevator")
						nextFloor <- orderArray[0][0]

						fmt.Println("Next floor is: ", orderArray[0][0])
					}
				} else {
					fmt.Println("I did not have lowest cost and did not take the order")
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func Clear_orders(orderFinished chan bool, nextFloor chan int, message chan string) {
	//send kvittering for ekstern ordre på nettverket
	//motta kvittering
	//fjern ordre fra kø

	Button := BUTTON_COMMAND
	for {
		if numberofOrders > 0 { // go to floor and remove order from que when finished

			select {
			case orderFinished_i := <-orderFinished:
				if orderFinished_i == true {
					if orderArray[1][0] == 0 { // Internal
						Button = BUTTON_CALL_UP
						//clear external order array
						for i := 0; i < ext_numberofOrders; i++{
							if ext_orderArray[1][i].number == 0{
								if ext_orderArray[0][i].number == CurrentFloor{
									message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[0][i].number), 10), strconv.FormatInt(int64(ext_orderArray[1][i].number), 10)}, ",")
								}
							}
						}
					} else if orderArray[1][0] == 1 {
						Button = BUTTON_CALL_DOWN
						//clear external order array
						for i := 0; i < ext_numberofOrders; i++{
							if ext_orderArray[1][i].number == 1{
								if ext_orderArray[0][i].number == CurrentFloor{
									message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[0][i].number), 10), strconv.FormatInt(int64(ext_orderArray[1][i].number), 10)}, ",")
								}
							}
						}
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

func Resend_externalorders(message chan string){
	for{
		now := time.Now()
		if ext_numberofOrders > 0{
			for i := 0; i < ext_numberofOrders; i++{
			//if external order has timed out, issue a new auction
				if now.Sub(ext_orderArray[0][i].starttime) > 20000000000 {
					fmt.Println("Order times out, sending new order to: ", ext_orderArray[0][i])
					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(ext_orderArray[0][i].number), 10), strconv.FormatInt(int64(ext_orderArray[1][i].number), 10)}, ",")
					ext_orderArray[0][i].starttime = time.Now()
				}
			}
		}
	time.Sleep(time.Second * 2)
	}
}
func main() {

	nextFloor := make(chan int, 10)
	orderFinished := make(chan bool, 1)
	up_button := make(chan int, 10)
	down_button := make(chan int, 10)
	internal_button := make(chan int, 10)
	message := make(chan string, 20)
	recievedmessage := make(chan string, 20)

	runtime.GOMAXPROCS(runtime.NumCPU())

	Init_system()
	fmt.Println("Init finished")


	go Broadcast(message, recievedmessage)
	go Elevator_driver(nextFloor, orderFinished)
	go TCP_sender(message, recievedmessage)
	go Local_orders(internal_button, nextFloor, orderFinished)
	go Handle_buttons(up_button, down_button, internal_button)
	go External_orders(message, up_button, down_button, nextFloor)
	go Displayfloor()
	go TCP_listener(recievedmessage)
	go Incomming_message(recievedmessage, message)
	go Assess_cost(nextFloor)
	go Clear_orders(orderFinished, nextFloor, message)
	go Resend_externalorders(message)


	deadChan := make(chan bool, 1)
	<-deadChan
}
