package Heismodul

import (
	. "./driver"
	"fmt"
	"os"
	"time"
)

//const N_FLOORS int = 4

func Displayfloor() {
	var oldFloor int = -1
	var newFloor int
	var floor int
	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			newFloor = floor
			if newFloor != oldFloor {
				oldFloor = newFloor
				//fmt.Println(newFloor + 1)
				Elev_set_floor_indicator(newFloor)
			}
		}
	}
}

func Init_floor() int {
	Elev_init()
	var oldFloor int = -1
	var newFloor int
	var floor int
	var foundFloor bool = false
	Elev_set_motor_direction(DIRN_UP)
	startTime := time.Now()
	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			newFloor = floor
			if newFloor != oldFloor {
				oldFloor = newFloor
				foundFloor = true
				Elev_set_motor_direction(DIRN_STOP)
				fmt.Println("FoundFloor")
				break
			}
		}
		currentTime := time.Now()
		if 2000000000 <= currentTime.Sub(startTime) {
			break
		}
	}

	if foundFloor == false {

		Elev_set_motor_direction(DIRN_DOWN)
		startTime = time.Now()
		for {
			floor = Elev_get_floor_sensor_signal()
			if floor >= 0 {
				newFloor = floor
				if newFloor != oldFloor {
					oldFloor = newFloor
					foundFloor = true
					Elev_set_motor_direction(DIRN_STOP)
					fmt.Println("FoundFloor")
					break
				}
			}
			currentTime := time.Now()
			if 2000000000 <= currentTime.Sub(startTime) {
				fmt.Println("FAILURE, move elevator away from endstops!")
				os.Exit(1)
			}
		}
	}
	Elev_set_motor_direction(DIRN_STOP)

	return int(oldFloor)
}

/*
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
*/
func Kjør_heis(nextFloor chan int, orderFinished chan bool) {
	var State int = 0
	var Finished = false
	var lastFloor int
	for {
		select {
		case nextFloor_i := <-nextFloor:
			Finished = false
			for Finished == false {
				var currentFloor = Elev_get_floor_sensor_signal()
				if currentFloor >= 0 {
					lastFloor = currentFloor
				}
				if State == 0 {

					if (lastFloor < nextFloor_i) && (nextFloor_i <= 3) {
						Elev_set_motor_direction(DIRN_UP)

					} else if (lastFloor > nextFloor_i) && (nextFloor_i >= 0) {
						Elev_set_motor_direction(DIRN_DOWN)

					} else {
						Elev_set_motor_direction(DIRN_STOP)
					}

					if currentFloor == nextFloor_i {
						Elev_set_motor_direction(DIRN_STOP)
						State = 1
						fmt.Println("State 1")
					}

				} else if State == 1 {
					Elev_set_door_open_lamp(1)
					orderFinished <- true
					time.Sleep(time.Second * 1)
					Elev_set_door_open_lamp(0)
					fmt.Println("Ready for new floor")
					State = 0
					Finished = true
				} else {
					Elev_set_motor_direction(DIRN_STOP)
				}
			}
		default:
			Elev_set_motor_direction(DIRN_STOP)
		}
	}
}

/*
func main() {
	//Elev_init()
	nextFloor := make(chan int, 1)
	orderFinished := make(chan bool, 1)

	go Intern_ordre(nextFloor, orderFinished)
	//go Kjør_heis(nextFloor, orderFinished)
	//go Displayfloor()
	deadChan := make(chan bool, 1)
	<-deadChan
}
*/
