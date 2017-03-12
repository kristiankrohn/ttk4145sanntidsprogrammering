package Elevatormodule

import (
	. "./driver"
	"fmt"
	"time"
)

const N_FLOORS int = 4

var CurrentFloor int

func Display_floor() {
	var oldFloor int = -1
	var newFloor int
	var floor int
	for {
		floor = Elev_get_floor_sensor_signal()
		if floor >= 0 {
			newFloor = floor
			if newFloor != oldFloor {
				oldFloor = newFloor
				Elev_set_floor_indicator(newFloor)
				CurrentFloor = newFloor
			}
		}
		time.Sleep(time.Millisecond * 10)

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

					break
				}
			}
			currentTime := time.Now()
			if 2000000000 <= currentTime.Sub(startTime) {
				fmt.Println("FAILURE, move elevator away from endstops!")
				time.Sleep(time.Second * 1)
			}
		}
	}
	Elev_set_motor_direction(DIRN_STOP)
}

func Handle_buttons(upButton chan int, downButton chan int, internalButton chan int) {

	var buttonPress [3][N_FLOORS]int
	var buttonRelease [3][N_FLOORS]int

	for i := 0; i < 3; i++ {
		for j := 0; j < N_FLOORS; j++ {
			buttonRelease[i][j] = 0
		}
	}
	button := BUTTON_CALL_UP
	for {
		for dir := 0; dir <= 2; dir++ {
			if dir == 0 {
				button = BUTTON_CALL_UP
			} else if dir == 1 {
				button = BUTTON_CALL_DOWN
			} else {
				button = BUTTON_COMMAND
			}

			for floor := 0; floor < N_FLOORS; floor++ { // read buttonpress and put on channel
				buttonPress[dir][floor] = Elev_get_button_signal(button, floor) //UP == 0, DOWN == 1
				if (buttonPress[dir][floor] == 1) && (buttonRelease[dir][floor] == 0) {
					buttonRelease[dir][floor] = 1

					if dir == 0 {

						upButton <- floor

					} else if dir == 1 {

						downButton <- floor

					} else {
						if floor != CurrentFloor {
							internalButton <- floor
						} else if Elev_get_floor_sensor_signal() != -1 {
							internalButton <- floor

						}
					}

				} else if (buttonPress[dir][floor] == 0) && (buttonRelease[dir][floor] == 1) {

					buttonRelease[dir][floor] = 0
				}
			}
		}
		time.Sleep(time.Millisecond * 10)

	}
}

func Elevator_driver(nextFloor chan int, orderFinished chan bool, stopElevator chan bool) {
	//this is the module acually controlling the elevator
	var State int = 0
	var Finished = false
	var nextFloor_i int = -1
	var stopElevator_i bool = false

	for {

		select {
		case nextFloor_ci := <-nextFloor:
			{
				nextFloor_i = int(nextFloor_ci)
				Finished = false

				if (CurrentFloor < nextFloor_i) && (nextFloor_i <= 3) {
					Elev_set_motor_direction(DIRN_UP)

				} else if (CurrentFloor > nextFloor_i) && (nextFloor_i >= 0) {
					Elev_set_motor_direction(DIRN_DOWN)

				} else {
					Elev_set_motor_direction(DIRN_STOP)

				}
			}
		case stopElevator_ci := <-stopElevator:
			{
				stopElevator_i = bool(stopElevator_ci)
				stopElevator_i = stopElevator_i
				fmt.Println("Stopping elevator on the fly")

				for stopElevator_i == true {
					Elev_set_motor_direction(DIRN_STOP)
					Elev_set_door_open_lamp(1)

					time.Sleep(time.Second * 2)
					Elev_set_door_open_lamp(0)


					stopElevator_i = false
				}
				if stopElevator_i == false {
					if (CurrentFloor < nextFloor_i) && (nextFloor_i <= 3) {
						Elev_set_motor_direction(DIRN_UP)

					} else if (CurrentFloor > nextFloor_i) && (nextFloor_i >= 0) {
						Elev_set_motor_direction(DIRN_DOWN)

					} else {
						Elev_set_motor_direction(DIRN_STOP)

					}
				}
			}

		default: // intentionally placed default, we're polling this one
			if State == 0 {

				if (CurrentFloor == nextFloor_i) && (Finished == false) {

					Elev_set_motor_direction(DIRN_STOP)
					State = 1
				}
			} else if State == 1 {
				Elev_set_door_open_lamp(1)

				State = 2

			} else if State == 2 {
				time.Sleep(time.Second * 2)
				State = 3

			} else if State == 3 {
				orderFinished <- true
				Elev_set_door_open_lamp(0)
				State = 0
				Finished = true
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}
