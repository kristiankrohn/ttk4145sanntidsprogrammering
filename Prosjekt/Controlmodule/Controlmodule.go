package Controlmodule

import (
	. "./Elevatormodule"
	. "./Elevatormodule/driver"
	. "./Networkmodule"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//const N_FLOORS int = 4 Define this in Heismodul.go
//const numberofelevators int = 10 // Declared in networkmodule
const arraysize int = N_FLOORS * 10 // number of buttons as a function of number of floors plus a litte margin

var orderArray [arraysize]Orderentry
var ext_orderArray [arraysize]Extentry
var cost_array [arraysize][Numberofelevators]Costentry
var numberofCosts [arraysize]Costnumber
var numberofOrders int
var ext_numberofOrders int

//var currentFloor int - Declared in elevatormodule

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
		}

		dataFile.Close()
		Button := BUTTON_CALL_UP
		for i := 0; (orderArray[i].Floor > -1) || (i == (arraysize - 1)); i++ {
			numberofOrders = i + 1

			if orderArray[i].Button == 0 {
				Button = BUTTON_CALL_UP
			} else if orderArray[i].Button == 1 {
				Button = BUTTON_CALL_DOWN
			} else {
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
}

func Calculate_cost(floor int, calldirection int) int {
	//very simple costfunction based on direction of travel, distance and number of orders
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
	//read buttonpress from channel
	//check if order is already in que
	// add order to que if its not already in
	var orderMatch bool = false

	for {
		select {
		case internal_call := <-internal_button:
			{

				orderMatch = false
				newOrder := int(internal_call)
				for j := 0; j <= numberofOrders; j++ {
					if (orderArray[j].Floor == newOrder) && (orderArray[j].Button == 2) {
						orderMatch = true
						fmt.Println("Order already exist")
					}
				}

				if orderMatch == false {

					orderArray[numberofOrders] = Orderentry{newOrder, 2}
					fmt.Println("New order at floor: ", newOrder)
					Elev_set_button_lamp(BUTTON_COMMAND, orderArray[numberofOrders].Floor, 1)
					numberofOrders++

					Backup_localorders()

					if numberofOrders == 1 {
						nextFloor <- orderArray[0].Floor
						fmt.Println("Next floor is: ", orderArray[0].Floor)
					}
				}
			}
		default:
		}
		time.Sleep(time.Millisecond * 10)

	}
}

func External_orders(message chan string, up_button chan int, down_button chan int, nextFloor chan int) {
	//read buttonpress from channel
	//check if order is already in que
	//share order on network

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

					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10),
						strconv.FormatInt(int64(DIR), 10)}, ",")
					ext_orderArray[ext_numberofOrders] = Extentry{newOrder, DIR, time.Now()}
					ext_numberofOrders++
					Elev_set_button_lamp(BUTTON_CALL_UP, newOrder, 1)
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

					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(newOrder), 10),
						strconv.FormatInt(int64(DIR), 10)}, ",")
					ext_orderArray[ext_numberofOrders] = Extentry{newOrder, DIR, time.Now()}
					ext_numberofOrders++
					Elev_set_button_lamp(BUTTON_CALL_DOWN, newOrder, 1)
				}
			}
		}
	}
	time.Sleep(time.Millisecond * 10)

}

func Message_handler(recievedmessage chan string, message chan string) {
	//Sort incomiing messages based on messagecode and send out proper response

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

					fmt.Println("Recieved external order")

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
					message <- strings.Join([]string{strconv.FormatInt(int64(1), 10), strconv.FormatInt(int64(cost), 10),
						strconv.FormatInt(int64(button), 10)}, ",")

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
					if (Floor >= 0)||(Floor > 3){
						fmt.Println("Clear external button at floor: ", Floor)
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
					} else {
						fmt.Println("Invalid clear message")
						fmt.Println(newOrder)
						os.Exit(0)
					}
					if ext_numberofOrders < 0 {
						ext_numberofOrders = 0
					} // just to be sure
				}
			}
		default:
			time.Sleep(time.Millisecond * 10)

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
				min := cost_array[i][0]
				fmt.Println("Order auction ended, assessing cost")
				for j := 0; j <= numberofCosts[i].number; j++ {

					//Sjekker hvilket bidrag som har lavest kost
					if cost_array[i][j].cost < min.cost {
						min = cost_array[i][j]
						//Dersom kosten er den samme som vår kost, men vår IP er lavere, så tar vi oppdraget
					} else if min.cost == cost_array[i][j].cost {
						fmt.Println("myIP = ", myIP)
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
					fmt.Println("I did not have lowest cost or cost&IP and did not take the order")
				}
				numberofCosts[i].number = 0
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func Clear_orders(orderFinished chan bool, nextFloor chan int, message chan string, stopElevator chan bool) {
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
									message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10),
										strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
								}
							}
						}
					} else if orderArray[0].Button == 1 {
						Button = BUTTON_CALL_DOWN
						//clear external order array
						for i := 0; i < ext_numberofOrders; i++ {
							if ext_orderArray[i].Button == 1 {
								if ext_orderArray[i].Floor == CurrentFloor {
									message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10),
										strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
								}
							}
						}
					} else {
						Button = BUTTON_COMMAND
					}
					fmt.Println("Order to floor: ", orderArray[0].Floor, " finished, removed from que")
					Elev_set_button_lamp(Button, orderArray[0].Floor, 0)
					
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
				{

					//on the fly check for orders that can be completed
					if numberofOrders > 1 {
						for i := 1; i < numberofOrders; i++ {
							// 	direction up = 0
							//	direction down = 1
							//	direction internal = 2
							//
							var direction int
							if orderArray[0].Floor > CurrentFloor { // going uo
								direction = 0
							} else if orderArray[0].Floor < CurrentFloor {
								direction = 1
							}
							direction = direction
							if orderArray[i].Floor == CurrentFloor {
								if (orderArray[i].Button == direction) || (orderArray[i].Button == 2) {
									stopElevator <- true

									fmt.Println("Current floor:", CurrentFloor)

									Button := BUTTON_COMMAND
									if orderArray[i].Button == 0 {
										Button = BUTTON_CALL_UP
										message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10),
											strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
									} else if orderArray[i].Button == 1 {
										Button = BUTTON_CALL_DOWN
										message <- strings.Join([]string{strconv.FormatInt(int64(2), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10),
											strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
									} else {
										Button = BUTTON_COMMAND
									}
									fmt.Println("Button to close: ", orderArray[i].Button)
									Elev_set_button_lamp(Button, CurrentFloor, 0)
									//Remove order from array
									for j := i; j < numberofOrders; j++ {
										orderArray[j] = orderArray[j+1]
									}
									numberofOrders--
									Backup_localorders()

									fmt.Println("Cleared order on the fly, remaining orders: ", numberofOrders)
								}
							}
						}
					}
				}
			}
		}
		time.Sleep(time.Millisecond * 10)

	}
}

func Resend_externalorders(message chan string) {
	//if external order has timed out, issue a new auction
	for {
		now := time.Now()
		if ext_numberofOrders > 0 {
			for i := 0; i < ext_numberofOrders; i++ {

				if now.Sub(ext_orderArray[i].starttime) > 30000000000 {
					fmt.Println("Order times out, sending new order to: ", ext_orderArray[i].Floor)
					message <- strings.Join([]string{strconv.FormatInt(int64(0), 10), strconv.FormatInt(int64(ext_orderArray[i].Floor), 10),
						strconv.FormatInt(int64(ext_orderArray[i].Button), 10)}, ",")
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
	} else {

		// serialize the data
		dataEncoder := gob.NewEncoder(dataFile)

		err = dataEncoder.Encode(&orderArray)
		if err != nil {
			fmt.Println(err)
		}
		dataFile.Close()
	}
}
