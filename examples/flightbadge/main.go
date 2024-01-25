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

var (
	drone *tello.Tello

	droneconnected bool
	takeoff        bool
	direction      int
)

const speed = 30

var logo = `
  ___ _ _      _   _      
 | __| (_)__ _| |_| |_    
 | _|| | / _\ | ' \  _|   
 |_|_|_|_\__, |_||_\__|   
 | _ ) __|___/| |__ _ ___ 
 | _ \/ _\ / _\ / _\ / -_)
 |___/\__,_\__,_\__, \___|
                |___/     
`

func main() {
	setupDisplay()
	time.Sleep(3 * time.Second)

	terminalOutput("enable wireless adapter...")

	link, _ := probe.Probe()

	err := link.NetConnect(&netlink.ConnectParams{
		Ssid:       ssid,
		Passphrase: pass,
	})

	if err != nil {
		failMessage(err.Error())
	}

	drone = tello.New("8888")
	connectDrone()

	go readControls()
	controlDrone()
}

func connectDrone() {
	terminalOutput("Starting drone...")
	if err := drone.Start(); err != nil {
		failMessage(err.Error())
	}

	terminalOutput("Drone started.")

	time.Sleep(1 * time.Second)

	// terminalOutput("Starting video...")
	// if err := drone.StartVideo(); err != nil {
	// 	failMessage(err.Error())
	// }
	// terminalOutput("Video started.")

	droneconnected = true
}

func controlDrone() {
	for {
		if !droneconnected {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		switch direction {
		case directionForward:
			drone.Forward(speed)
		case directionBackward:
			drone.Backward(speed)
		default:
			drone.Forward(0)
		}

		switch direction {
		case directionLeft:
			drone.Left(speed)
		case directionRight:
			drone.Right(speed)
		default:
			drone.Right(0)
		}

		switch direction {
		case directionUp:
			drone.Up(speed)
		case directionDown:
			drone.Down(speed)
		default:
			drone.Up(0)
		}

		switch direction {
		case directionTurnLeft:
			drone.CounterClockwise(speed)
		case directionTurnRight:
			drone.Clockwise(speed)
		default:
			drone.Clockwise(0)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func failMessage(msg string) {
	for {
		terminalOutput(msg)
		time.Sleep(1 * time.Second)
	}
}
