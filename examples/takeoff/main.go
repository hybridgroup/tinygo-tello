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

	connectToAP(droneConnected)
}

func droneConnected() {
	println("Starting drone")
	drone.Start()

	time.Sleep(5 * time.Second)

	println("Taking off")
	drone.TakeOff()
	time.Sleep(5 * time.Second)

	println("Landing")
	drone.Land()
}

// connect to drone wifi
func connectToAP(connectHandler func()) {
	var err error
	time.Sleep(2 * time.Second)
	for i := 0; i < 3; i++ {
		println("Connecting to " + ssid)
		err = adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
		if err != nil {
			println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		// success
		println("Connected.")
		time.Sleep(3 * time.Second)
		if connectHandler != nil {
			connectHandler()
		}
	}

	// couldn't connect to AP
	failMessage(err.Error())
}

func failMessage(msg string) {
	for {
		println(msg)
		time.Sleep(1 * time.Second)
	}
}
