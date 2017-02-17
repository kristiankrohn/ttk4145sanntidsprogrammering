package main

import (
	. "./driver"
	"fmt"
	"os"
	"time"
)

const N_FLOORS int = 4

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

func TestElevator() {
	fmt.Println("Press STOP button to stop elevator and exit program")
	Elev_set_motor_direction(DIRN_UP)
	var floor int
	for {

		floor = Elev_get_floor_sensor_signal()

		if floor == (N_FLOORS - 1) {
			Elev_set_motor_direction(DIRN_DOWN)
		} else if floor == 0 {
			Elev_set_motor_direction(DIRN_UP)
		}

		// Stop elevator and exit program if the stop button is pressed
		if Elev_get_stop_signal() == 1 {
			Elev_set_motor_direction(DIRN_STOP)
			break
		}
	}
}

func Init_floor() int {
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
				break
			}
		}
		currentTime := time.Now()
		if 1500000000 <= currentTime.Sub(startTime) {
			break
		}
	}

	if foundFloor == false {
		Elev_init()
		Elev_set_motor_direction(DIRN_DOWN)
		startTime = time.Now()
		for {
			floor = Elev_get_floor_sensor_signal()
			if floor >= 0 {
				newFloor = floor
				if newFloor != oldFloor {
					oldFloor = newFloor
					foundFloor = true
					break
				}
			}
			currentTime := time.Now()
			if 1500000000 <= currentTime.Sub(startTime) {
				fmt.Println("FAILURE, move elevator away from endstops!")
				os.Exit(1)
			}
		}
	}

	Elev_set_motor_direction(DIRN_STOP)
	return int(oldFloor)

}

func InternalTest() {
	go Displayfloor()
	var currentFloor = Init_floor()
	var floor int

	var numberofOrders int = 0
	var orderArray [N_FLOORS + 1]int
	for j := 0; j <= N_FLOORS; j++ {
		orderArray[j] = -1
	}
	var newOrder int
	var orderMatch bool
	var orderFinished bool
	var oldOrderFinished bool = false

	var buttonPress [4]int
	buttonRelease := [4]int{0, 0, 0, 0}

	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			currentFloor = floor
		}

		if numberofOrders > 0 { // go to floor and remove order from que when finished
			orderFinished = GoToFloor(currentFloor, orderArray[0])
			if (orderFinished == true) && (oldOrderFinished == false) {
				oldOrderFinished = true
				Elev_set_button_lamp(BUTTON_COMMAND, orderArray[0], 0)
				fmt.Println("Order to floor: ", orderArray[0], " finished, removed from que")
				for i := 0; i < numberofOrders; i++ {
					orderArray[i] = orderArray[i+1]
				}

				numberofOrders--
				fmt.Println("Number of order: ", numberofOrders)
			} else if (orderFinished == false) && (oldOrderFinished == true) {
				oldOrderFinished = false
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
				}

			} else if (buttonPress[i] == 0) && (buttonRelease[i] == 1) {
				//fmt.Println("New buttonrelease at: ", i)
				buttonRelease[i] = 0
			}
		}
	}
}

func GoToFloor(currentFloor int, nextFloor int) bool {
	var startTime time.Time
	if currentFloor < nextFloor {
		Elev_set_motor_direction(DIRN_UP)
		startTime = time.Now()
	} else if currentFloor > nextFloor {
		Elev_set_motor_direction(DIRN_DOWN)
		startTime = time.Now()
	} else {
		Elev_set_motor_direction(DIRN_STOP)
	}

	var floor = Elev_get_floor_sensor_signal()
	if floor == nextFloor {
		Elev_set_motor_direction(DIRN_STOP)
		Elev_set_door_open_lamp(1)
		fmt.Println(startTime)
		currentTime := time.Now()
		if 1000000000 <= currentTime.Sub(startTime) {
			Elev_set_door_open_lamp(0)
			fmt.Println("Ready for new floor")
			return true
		} else {return false}
	}
	return false
}

func main() {
	Elev_init()

	InternalTest()
	//TestElevator()

	//deadChan := make(chan bool, 1)
	//<-deadChan
}
