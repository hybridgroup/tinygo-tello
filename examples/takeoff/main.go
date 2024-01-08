package main

import (
	"log"
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

func main() {
	link, _ := probe.Probe()

	err := link.NetConnect(&netlink.ConnectParams{
		Ssid:       ssid,
		Passphrase: pass,
	})
	if err != nil {
		log.Fatal(err)
	}

	drone := tello.New("8888")
	println("Starting drone")
	if err := drone.Start(); err != nil {
		for {
			println(err.Error())
			time.Sleep(1 * time.Second)
		}
	}

	time.Sleep(5 * time.Second)

	println("Taking off")
	if err := drone.TakeOff(); err != nil {
		for {
			println(err.Error())
			time.Sleep(1 * time.Second)
		}
	}

	time.Sleep(5 * time.Second)

	println("Landing")
	drone.Land()
}
