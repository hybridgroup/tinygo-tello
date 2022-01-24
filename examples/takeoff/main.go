package main

import (
	"time"

	tello "github.com/hybridgroup/tinygo-tello"
)

// Tello drone info here
const ssid = "TELLO-C48E59"
const pass = ""

var (
	drone *tello.Tello
)

func main() {
	a := initAdaptor()
	drone = tello.New(a, "8888")

	connectToAP(connectDrone)
}

func connectDrone() {
	println("Starting drone")
	drone.Start()

	time.Sleep(5 * time.Second)

	println("Taking off")
	drone.TakeOff()
	time.Sleep(5 * time.Second)

	println("Landing")
	drone.Land()
}

// connect to access point
func connectToAP(connectHandler func()) {
	time.Sleep(2 * time.Second)
	println("Connecting to " + ssid)
	err := adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
	if err != nil { // error connecting to AP
		for {
			println(err)
			time.Sleep(1 * time.Second)
		}
	}

	println("Connected.")

	time.Sleep(3 * time.Second)
	if connectHandler != nil {
		connectHandler()
	}
}

func message(msg string) {
	println(msg, "\r")
}
