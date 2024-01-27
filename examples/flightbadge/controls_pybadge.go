//go:build pybadge

package main

import (
	"time"

	tello "github.com/hybridgroup/tinygo-tello"
	"tinygo.org/x/drivers/shifter"
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

var (
	shifted bool
	flip    bool
)

func readControls() {
	buttons := shifter.NewButtons()
	buttons.Configure()

	for {
		buttons.ReadInput()

		// takeoff
		if buttons.Pins[shifter.BUTTON_START].Get() {
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
		if buttons.Pins[shifter.BUTTON_B].Get() {
			terminalOutput("landing")
			err := drone.Land()
			if err != nil {
				terminalOutput(err.Error())
			}
			takeoff = false
		}

		// front flip
		if buttons.Pins[shifter.BUTTON_SELECT].Get() {
			if !flip {
				terminalOutput("flip")
				err := drone.Flip(tello.FlipFront)
				if err != nil {
					terminalOutput(err.Error())
				}
				flip = true
			}
		}

		// hold down button A to shift to access second set of arrow commands
		if buttons.Pins[shifter.BUTTON_A].Get() {
			shifted = true

			// reset flip
			flip = false
		} else {
			shifted = false
		}

		direction = directionNone

		if buttons.Pins[shifter.BUTTON_LEFT].Get() {
			if shifted {
				direction = directionTurnLeft
			} else {
				direction = directionLeft
			}
		}

		if buttons.Pins[shifter.BUTTON_UP].Get() {
			if shifted {
				direction = directionUp
			} else {
				direction = directionForward
			}
		}

		if buttons.Pins[shifter.BUTTON_DOWN].Get() {
			if shifted {
				direction = directionDown
			} else {
				direction = directionBackward
			}
		}

		if buttons.Pins[shifter.BUTTON_RIGHT].Get() {
			if shifted {
				direction = directionTurnRight
			} else {
				direction = directionRight
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}
