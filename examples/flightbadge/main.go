// TinyGo flight control for the Tello drone in the form of a badge.
package main

import (
	"time"

	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/netlink/probe"

	tello "github.com/hybridgroup/tinygo-tello"
)

// Tello drone info here
var (
	ssid string
	pass string
)

const (
	speed   = 30
	center  = 660
	detente = 300
)

var (
	drone *tello.Tello
)

func main() {
	println("Connecting to drone...")

	link, _ := probe.Probe()

	err := link.NetConnect(&netlink.ConnectParams{
		Ssid:       ssid,
		Passphrase: pass,
	})

	if err != nil {
		failMessage(err.Error())
	}

	drone = tello.New("8888")

	initDisplay()
	go handleDisplay()

	initControls()
	go readControls()

	connectDrone()
	controlDrone()
}

func connectDrone() {
	println("Starting drone...")
	if err := drone.Start(); err != nil {
		failMessage(err.Error())
	}

	println("Drone started.")

	time.Sleep(1 * time.Second)

	println("Starting video...")
	if err := drone.StartVideo(); err != nil {
		failMessage(err.Error())
	}
	println("Video started.")

	droneconnected = true
}

func controlDrone() {
	startvid := true

	for {
		switch {
		case b1push:
			println("takeoff")
			err := drone.TakeOff()
			if err != nil {
				println(err)
			}

		case b2push:
			println("land")
			err := drone.Land()
			if err != nil {
				println(err)
			}
		}

		rightStick := getRightStick()
		switch {
		case rightStick.y+detente < center:
			drone.Backward(speed)
		case rightStick.y-detente > center:
			drone.Forward(speed)
		default:
			drone.Forward(0)
		}

		switch {
		case rightStick.x-detente > center:
			drone.Right(speed)
		case rightStick.x+detente < center:
			drone.Left(speed)
		default:
			drone.Right(0)
		}

		leftStick := getLeftStick()
		switch {
		case leftStick.y+detente < center:
			drone.Down(speed)
		case leftStick.y-detente > center:
			drone.Up(speed)
		default:
			drone.Up(0)
		}

		switch {
		case leftStick.x-detente > center:
			drone.Clockwise(speed)
		case leftStick.x+detente < center:
			drone.CounterClockwise(speed)
		default:
			drone.Clockwise(0)
		}

		if startvid {
			drone.StartVideo()
			startvid = false
		} else {
			startvid = true
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func failMessage(msg string) {
	failure = msg
	for {
		println(msg)
		time.Sleep(1 * time.Second)
	}
}
