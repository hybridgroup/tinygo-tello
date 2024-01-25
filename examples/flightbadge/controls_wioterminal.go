//go:build wioterminal

package main

import (
	"machine"

	"time"
)

const (
	directionNone = iota
	directionForward
	directionBackward
	directionLeft
	directionRight
	directionUp
	directionDown
	directionTurnLeft
	directionTurnRight
)

var shifted bool

func readControls() {
	machine.WIO_5S_UP.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_DOWN.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_LEFT.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_RIGHT.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.WIO_5S_PRESS.Configure(machine.PinConfig{Mode: machine.PinInput})

	machine.BUTTON_1.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.BUTTON_2.Configure(machine.PinConfig{Mode: machine.PinInput})
	machine.BUTTON_3.Configure(machine.PinConfig{Mode: machine.PinInput})

	for {
		// takeoff
		if machine.BUTTON_2.Get() {
			if !takeoff {
				terminalOutput("takeoff")
				err := drone.TakeOff()
				if err != nil {
					terminalOutput(err.Error())
				}
				takeoff = true
			}
		}

		// land
		if machine.BUTTON_3.Get() {
			terminalOutput("landing")
			err := drone.Land()
			if err != nil {
				terminalOutput(err.Error())
			}
			takeoff = false
		}

		// front flip
		if machine.WIO_5S_PRESS.Get() {
			//handleKey("t")
		}

		// hold down button A to shift to access second set of arrow commands
		if machine.BUTTON_2.Get() {
			shifted = true
		} else {
			shifted = false
		}

		direction = directionNone

		if machine.WIO_5S_LEFT.Get() {
			if shifted {
				direction = directionTurnLeft
			} else {
				direction = directionLeft
			}
		}

		if machine.WIO_5S_UP.Get() {
			if shifted {
				direction = directionUp
			} else {
				direction = directionForward
			}
		}

		if machine.WIO_5S_DOWN.Get() {
			if shifted {
				direction = directionDown
			} else {
				direction = directionBackward
			}
		}

		if machine.WIO_5S_RIGHT.Get() {
			if shifted {
				direction = directionTurnRight
			} else {
				direction = directionRight
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}
