package Heismodul

import (
	. "./driver"
	"fmt"
	"os"
	"time"
)

const N_FLOORS int = 4

func Displayfloor(current_floor_external chan int, current_floor_internal chan int) {
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
				current_floor_external <- newFloor
				current_floor_internal <- newFloor
				fmt.Println("Current floor is: ", oldFloor)
			}
		}
	}
}

func Init_elevator() {
	Elev_init()
	var oldFloor int = -1
	var newFloor int
	var floor int
	var foundFloor bool = false
	Elev_set_motor_direction(DIRN_UP)
	startTime := time.Now()
	fmt.Println("Looking for floor")
	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			newFloor = floor
			if newFloor != oldFloor {
				oldFloor = newFloor
				foundFloor = true
				Elev_set_motor_direction(DIRN_STOP)
				fmt.Println("FoundFloorUP")
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
					fmt.Println("FoundFloorDOWN")
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

	//return int(oldFloor)
}

func Current_floor(current_floor_external chan int, current_floor_internal chan int) {
	var floor int
	var newFloor int
	var oldFloor int
	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			newFloor = floor
			if newFloor != oldFloor {
				oldFloor = newFloor
				select {
				case current_floor_internal <- oldFloor:
				default:
				}
				select {
				case current_floor_external <- oldFloor:
				default:
				}
				fmt.Println("Current floor is: ", oldFloor)
			}
		}
	}
}

func Handle_buttons(up_button chan int, down_button chan int, internal_button chan int) {

	var buttonPress [3][N_FLOORS]int
	var buttonRelease [3][N_FLOORS]int

	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			buttonRelease[i][j] = 0
		}
	}
	Button := BUTTON_CALL_UP
	for {
		for DIR := 0; DIR <= 2; DIR++ {
			if DIR == 0 {
				Button = BUTTON_CALL_UP
			} else if DIR == 1 {
				Button = BUTTON_CALL_DOWN
			} else {
				Button = BUTTON_COMMAND
			}

			for FLOOR := 0; FLOOR < N_FLOORS; FLOOR++ { // read buttonpress and add order to que
				buttonPress[DIR][FLOOR] = Elev_get_button_signal(Button, FLOOR) //UP == 0, DOWN == 1
				if (buttonPress[DIR][FLOOR] == 1) && (buttonRelease[DIR][FLOOR] == 0) {
					buttonRelease[DIR][FLOOR] = 1
					fmt.Println("New buttonpress at: ", FLOOR)
					if DIR == 0 {
						up_button <- FLOOR
					} else if DIR == 1 {
						down_button <- FLOOR
					} else {
						internal_button <- FLOOR
					}

				} else if (buttonPress[DIR][FLOOR] == 0) && (buttonRelease[DIR][FLOOR] == 1) {
					//fmt.Println("New buttonrelease at: ", i)
					buttonRelease[DIR][FLOOR] = 0
				}
			}
		}
	}
}

func Elevator_driver(nextFloor chan int, orderFinished chan bool) {

	var State int = 0
	var Finished = false
	var currentFloor int

	for {
		floor := Elev_get_floor_sensor_signal()
		if floor >= 0 {
			currentFloor = floor
			//fmt.Println(currentFloor)
		}

		select {
		case nextFloor_i := <-nextFloor:
			fmt.Println("Going for next floor", nextFloor_i)
			Finished = false

			if (currentFloor < nextFloor_i) && (nextFloor_i <= 3) {
				Elev_set_motor_direction(DIRN_UP)
				fmt.Println("UP")
			} else if (currentFloor > nextFloor_i) && (nextFloor_i >= 0) {
				Elev_set_motor_direction(DIRN_DOWN)
				fmt.Println("DOWN")
			} else {
				Elev_set_motor_direction(DIRN_STOP)
				fmt.Println("STOP")
			}

			for Finished == false {
				floor := Elev_get_floor_sensor_signal()
				if floor >= 0 {
					currentFloor = floor
					//fmt.Println(currentFloor)
				}

				if State == 0 {

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
					//Elev_set_motor_direction(DIRN_STOP)
				}
			}
		default:
			//Elev_set_motor_direction(DIRN_STOP)
		}
	}
}
