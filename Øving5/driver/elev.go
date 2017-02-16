package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"

type elev_motor_direction_t int

const (
	DIRN_DOWN elev_motor_direction_t = -1
	DIRN_STOP elev_motor_direction_t = 0
	DIRN_UP   elev_motor_direction_t = 1
)

type elev_button_type_t int

const (
	BUTTON_CALL_UP elev_button_type_t = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

func Elev_init() {
	C.elev_init()
}

func Elev_set_motor_direction(dirn elev_motor_direction_t) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func Elev_set_button_lamp(button elev_button_type_t, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Elev_set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Elev_set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func Elev_set_stop_lamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func Elev_get_button_signal(button elev_button_type_t, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Elev_get_floor_sensor_signal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Elev_get_stop_signal() int {
	return int(C.elev_get_stop_signal())
}

func Elev_get_obstruction_signal() int {
	return int(C.elev_get_obstruction_signal())
}
