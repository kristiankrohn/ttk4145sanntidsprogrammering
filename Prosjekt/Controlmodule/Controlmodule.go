package Controlmodule

import (
	. "./Elevatormodule"
	. "./Elevatormodule/driver"
	. "./Networkmodule"
	"fmt"
	"strconv"
	"strings"
	"time"
	"os"
	"encoding/gob"
)

//const _FLOORS int = 4 Define this in Heismodul.go
//const numberofelevators int = 5 //must be higher than maximum number of possible elevators, or it will cause bufferoverflow
const arraysize int = N_FLOORS * 10 // number of buttons as a function of number of elevators

var orderArray [arraysize]Orderentry
var ext_orderArray [arraysize]Extentry
var cost_array [arraysize][Numberofelevators]Costentry
var numberofCosts [arraysize]Costnumber
var numberofOrders int
var ext_numberofOrders int

//var currentFloor int - Declared in Heismodul.go
type Orderentry struct {
	Floor  int
	Button int
}

type Costentry struct {
	cost int
	IP   int
}

type Costnumber struct {
	number    int
	starttime time.Time
}

type Extentry struct {
	Floor     int
	Button    int
	starttime time.Time
}

/* 	direction up = 0
*	direction down = 1
*	direction internal = 2
 */

func Init_system(nextFloor chan int) {
	Init_elevator()

	dataFile, err := os.Open("OrderBackup.gob")

 	if err != nil {
 		fmt.Println(err)
 		fmt.Println("Initializing array")
 		for j := 0; j < arraysize; j++ {
			orderArray[j] = Orderentry{-1, -1}
		}
		numberofOrders = 0
 	} else {
 	
 		dataDecoder := gob.NewDecoder(dataFile)
 		err = dataDecoder.Decode(&orderArray)

 		if err != nil {
 			fmt.Println(err)
 			//os.Exit(1)
 		}

 		dataFile.Close()
 		//fmt.Println("Opened this array from file: ", orderArray)
 		Button := BUTTON_CALL_UP
 		for i := 0; (orderArray[i].Floor > -1) || (i == (arraysize - 1)); i++{
 			numberofOrders = i + 1

 			if orderArray[i].Button == 0{
 				Button = BUTTON_CALL_UP
 			} else if orderArray[i].Button == 1{
 				Button = BUTTON_CALL_DOWN
 			} else{
 				Button = BUTTON_COMMAND
 			}
 			Elev_set_button_lamp(Button, orderArray[i].Floor, 1)
 		}
 		fmt.Println("Number of orders: ", numberofOrders)

 		if numberofOrders >= 1 {
 			nextFloor <- orderArray[0].Floor
 		} 
	}


	for j := 0; j < arraysize; j++ {
		ext_orderArray[j] = Extentry{-1, -1, time.Now()}
	}

	for i := 0; i < arraysize; i++ {
		for j := 0; j < Numberofelevators; j++ {
			cost_array[i][j] = Costentry{0, 0}
		}
	}
	initCostnumber := Costnumber{0, time.Now()}
	for k := 0; k < arraysize; k++ {
		numberofCosts[k] = initCostnumber
	}
	
	//Init_elevator()

}

func Calculate_cost(floor int, calldirection int) int {
	//gjør beregninger - eksterne og interne ordre
	//det enkleste er å se på det totale antall ordre og etasjer den skal kjøre
	//del resultat på nettverk
	var cost int
	var direction int = 10

	if numberofOrders == 0 {
		if floor > CurrentFloor {
			cost = floor - CurrentFloor
		} else {
			cost = CurrentFloor - floor
		}
	} else {
		if floor > CurrentFloor {
			if orderArray[0].Floor > CurrentFloor {
				direction = 1
			} else {
				direction = 2
			}
			cost = direction * (floor - CurrentFloor)
		} else {
			if orderArray[0].Floor > CurrentFloor {
				direction = 2
			} else {
				direction = 1
			}
			cost = direction * (CurrentFloor - floor)
		}
	}
	return cost * (numberofOrders + 1)
}

func Local_orders(internal_button chan int, nextFloor chan int, orderFinished chan bool) {

	//go Displayfloor()
	//go Elevator_driver(nextFloor, orderFinished)
	//go Handle_Buttons(up_button, down_button, internal_button)
	var orderMatch bool = false

	for {

		select {
		case internal_call := <-internal_button:
			{
				//fmt.Println("Local_orders is alive and reading from channel!")
				orderMatch = false
				newOrder := int(internal_call)
				for j := 0; j <= numberofOrders; j++ {
					//fmt.Println("From loop")
					if (orderArray[j].Floor == newOrder) && (orderArray[j].Button == 2) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {
					//if (orderMatch == false) && (currentFloor != newOrder) {
					//fmt.Println("Current floor is: ", currentFloor)
					orderArray[numberofOrders] = Orderentry{newOrder, 2}
					//orderArray[1][numberofOrders] = 2
					fmt.Println("New order at floor: ", newOrder)

					Elev_set_button_lamp(BUTTON_COMMAND, orderArray[numberofOrders].Floor, 1)
					numberofOrders++

					Backup_localorders()

					if numberofOrders == 1 {
						nextFloor <- orderArray[0].Floor
						//fmt.Println("Next floor is: ", orderArray[0].floor)
					} /*else if numberofOrders > 1 { // tried out some array sorting
					/* 	direction up = 0
					*	direction down = 1
					*	direction internal = 2
					*/
					/*
						if orderArray[0].floor < currentFloor{

						} else if currentFloor < orderArray[0].floor{

						}
						//fmt.Println(orderArray)

						A := orderArray[0]
						B := orderArray[1]
						orderArray[0] = B
						orderArray[1] = A

						fmt.Println("Postsort: ", orderArray)
						nextFloor <- orderArray[0].floor
					}*/
				}
				//fmt.Println("Orderhandling finished")
				//fmt.Println(orderArray)
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
					if (ext_orderArray[j].Floor == newOrder) && (ext_orderArray[j].Button == DIR) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {

					fmt.Println("New order at floor: ", newOrder)

					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
				}
			}

		case call_down := <-down_button:
			{
				orderMatch = false
				DIR = 1
				newOrder = int(call_down)
				for j := 0; j <= ext_numberofOrders; j++ {
					if (ext_orderArray[j].Floor == newOrder) && (ext_orderArray[j].Button == DIR) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {

					fmt.Println("New order at floor: ", newOrder)

					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10), strconv.FormatInt(int64(DIR), 10)}, ",")
				}
			}
		}
	}
}

func Incomming_message(recievedmessage chan string, message chan string) {
	//check om ordre allerede er i køen
	//legg til hvis ikke
	//fjern intern ordre fra kø når ordren er fullført

	//sorry for using strings :(
	//its a complete mess, but it works

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

				messagecode, err := strconv.ParseInt(slice[0], 10, 64)
				CheckError(err)

				if messagecode == 0 { // Recieving a new order and sening out the cost
					Floor, err := strconv.ParseInt(slice[1], 10, 64)
					CheckError(err)
					newOrder := int(Floor)
					direction, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)

					DIR := int(direction)
					if ext_numberofOrders < 0 {
						ext_numberofOrders = 0
					}

					if DIR == 0 {

						Elev_set_button_lamp(BUTTON_CALL_UP, newOrder, 1) // Flyttes etterhvert?
					} else if DIR == 1 {

						Elev_set_button_lamp(BUTTON_CALL_DOWN, newOrder, 1) // Flyttes etterhvert?
					}

					orderMatch := false

					for j := 0; j <= ext_numberofOrders; j++ {
						if (ext_orderArray[j].Floor == newOrder) && (ext_orderArray[j].Button == DIR) {
							orderMatch = true
							fmt.Println("Extorder already exist")
						}
					}
					if orderMatch == false { // we have a new order that doesn't exist in the array
						ext_orderArray[ext_numberofOrders] = Extentry{newOrder, DIR, time.Now()}
						ext_numberofOrders++

						fmt.Println("Remaining external order: ", ext_numberofOrders)
						fmt.Println("New external order array", ext_orderArray[0].Floor)
					}

					cost := Calculate_cost(newOrder, DIR)

					button := newOrder + (DIR+1)*(DIR+1)

					message <- strings.Join([]string{strconv.FormatInt(int64(1), 10), strconv.FormatInt(int64(cost), 10), strconv.FormatInt(int64(button), 10)}, ",")
					fmt.Println("Returning Cost: ", cost)

				} else if messagecode == 1 { // Gather costs and put them in costarray
					// Slicing string and convert to useful datatypes
					cost_i64, err := strconv.ParseInt(slice[1], 10, 64) // converting via i64
					CheckError(err)
					cost := int(cost_i64)
					button_i64, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)
					button := int(button_i64)
					lastaddressbytestring_pp := strings.Split(slice[3], ":")
					lastaddressbytestring := strings.Split(lastaddressbytestring_pp[0], ".")
					lastaddressbyte_i64, err := strconv.ParseInt(lastaddressbytestring[3], 10, 64)
					CheckError(err)
					lastaddressbyte := int(lastaddressbyte_i64)
					if numberofCosts[button].number <= Numberofelevators {
						fmt.Println("Adding cost to array", cost, button, numberofCosts[button].number) //Cost_array = arraysize * numberofelevators

						// Adding cost to costarray and starting timer if it is first cost added
						cost_array[button][numberofCosts[button].number] = Costentry{cost, lastaddressbyte}
						if numberofCosts[button].number == 0 {
							numberofCosts[button].starttime = time.Now()
						}
						numberofCosts[button].number++
					} else {
						fmt.Println("Cost_array is full")
					}
				} else if messagecode == 2 { //Ext_order is completed and removed from ext_array
					Floor_i64, err := strconv.ParseInt(slice[1], 10, 64)
					CheckError(err)

					DIR_i64, err := strconv.ParseInt(slice[2], 10, 64)
					CheckError(err)

					Floor := int(Floor_i64)
					DIR := int(DIR_i64)

					Button := BUTTON_CALL_UP
					if DIR == 0 {
						Button = BUTTON_CALL_UP
					} else {
						Button = BUTTON_CALL_DOWN
					}

					Elev_set_button_lamp(Button, Floor, 0)

					for i := 0; i < ext_numberofOrders; i++ {
						if (ext_orderArray[i].Floor == Floor) && (ext_orderArray[i].Button == DIR) {
							for j := i; j < ext_numberofOrders; j++ {
								ext_orderArray[j] = ext_orderArray[j+1]
								i = j
							}
							ext_numberofOrders--
							fmt.Print("Removed completed externalorder, remaining: ", ext_numberofOrders)
							fmt.Println(" First is to floor: ", ext_orderArray[0].Floor)
						}
					}

					if ext_numberofOrders < 0 {
						ext_numberofOrders = 0
					} // just to be sure
				}
			}
		default:
		}
	}
}

func Assess_cost(nextFloor chan int) {
	//sammenlign innkommende resultat
	//Vurder om vi skal ta ordre og legge den til i intern ordrekø
	for {
		myIP := Last_byte_of_my_IP()
		for i := 0; i < arraysize; i++ {
			now := time.Now()
			//fmt.Println("checking for ordertimeouts: ", i)
			if (now.Sub(numberofCosts[i].starttime) > 1000000000) && (numberofCosts[i].number > 0) { // check for timeout, if timeout, assess costarray
				min := Costentry{9000, myIP}
				fmt.Println("Order auction ended, assessing cost")
				for j := 0; j < numberofCosts[i].number; j++ {

					//Sjekker hvilket bidrag som har lavest kost
					if cost_array[i][j].cost < min.cost {
						min = cost_array[i][j]
						//Dersom kosten er den samme som vår kost, men vår IP er lavere, så tar vi oppdraget
					} else if min.cost == cost_array[i][j].cost {
						fmt.Println("min.IP = ", min.IP)
						if myIP <= min.IP {
							min = cost_array[i][j]
							fmt.Println("Same cost, i have lowest IP of: ", myIP)
						} else {
							fmt.Println("Same cost but i have higer IP of: ", myIP)
						}
					}

				}
				
				if (myIP == min.IP) || (numberofCosts[i].number == 1) {

					fmt.Println("I have lowest cost or cost&IP and taking order, the number in que is", numberofOrders)

					if i < 4 {
						orderArray[numberofOrders] = Orderentry{(i - 1), 0}

					} else {
						orderArray[numberofOrders] = Orderentry{(i - 4), 1}

					}
					numberofOrders++

					Backup_localorders()

					if numberofOrders == 1 {
						fmt.Println("This is first order and sending to elevator")
						nextFloor <- orderArray[0].Floor

						fmt.Println("Next floor is: ", orderArray[0].Floor)
					}
				} else {
					fmt.Println("I did not have lowest cost and did not take the order")
				}
				numberofCosts[i].number = 0
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
				fmt.Println("Removing order")
				if orderFinished_i == true {
					if orderArray[0].Button == 0 { // Internal
						Button = BUTTON_CALL_UP
						//clear external order array
						for i := 0; i < ext_numberofOrders; i++ {
							if ext_orderArray[i].Button == 0 {
								if ext_orderArray[i].Floor == CurrentFloor {
									message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10), strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
								}
							}
						}
					} else if orderArray[0].Button == 1 {
						Button = BUTTON_CALL_DOWN
						//clear external order array
						for i := 0; i < ext_numberofOrders; i++ {
							if ext_orderArray[i].Button == 1 {
								if ext_orderArray[i].Floor == CurrentFloor {
									message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10), strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
								}
							}
						}
					} else {
						Button = BUTTON_COMMAND
					}
					Elev_set_button_lamp(Button, orderArray[0].Floor, 0)
					fmt.Println("Order to floor: ", orderArray[0].Floor, " finished, removed from que")
					for i := 0; i < numberofOrders; i++ {
						orderArray[i] = orderArray[i+1]
					}
					numberofOrders--
					fmt.Println("Number of orders: ", numberofOrders)
					//fmt.Println(orderArray)

					Backup_localorders()

					if numberofOrders >= 1 {
						nextFloor <- orderArray[0].Floor
						fmt.Println("Next floor is: ", orderArray[0].Floor)
					}

				}
			default:
			}
		}
	}
}

func Resend_externalorders(message chan string) {
	for {
		now := time.Now()
		if ext_numberofOrders > 0 {
			for i := 0; i < ext_numberofOrders; i++ {
				//if external order has timed out, issue a new auction
				if now.Sub(ext_orderArray[i].starttime) > 20000000000 {
					fmt.Println("Order times out, sending new order to: ", ext_orderArray[i].Floor)
					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10), strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
					ext_orderArray[i].starttime = time.Now()
				}
			}
		}
		time.Sleep(time.Second * 2)
	}
}

func Backup_localorders() {
	dataFile, err := os.Create("OrderBackup.gob")
 	if err != nil {
 		fmt.Println(err)
 		//os.Exit(1)
 	} else{

      	// serialize the data
 		dataEncoder := gob.NewEncoder(dataFile)

 		err = dataEncoder.Encode(&orderArray)
 		if err != nil {
 			fmt.Println(err)
 			//os.Exit(1)
 		}

 		dataFile.Close()
 	}
}
