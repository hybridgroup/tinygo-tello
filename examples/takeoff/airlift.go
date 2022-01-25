//go:build airlift
// +build airlift

package main

import (
	"machine"

	"tinygo.org/x/drivers/wifinina"
)

var (
	// default interface for the Arduino Nano33 IoT.
	spi = machine.SPI0

	// ESP32/ESP8266 chip that has the WIFININA firmware flashed on it
	adaptor *wifinina.Device
)

func initAdaptor() *wifinina.Device {
	// Configure SPI for 8Mhz, Mode 0, MSB First
	spi.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
		SDO:       machine.SPI0_SDO_PIN,
		SDI:       machine.SPI0_SDI_PIN,
		SCK:       machine.SPI0_SCK_PIN,
	})

	adaptor = wifinina.New(spi,
		machine.D13,
		machine.D11,
		machine.D10,
		machine.D12)
	adaptor.Configure()

	return adaptor
}
