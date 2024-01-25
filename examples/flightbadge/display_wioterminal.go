//go:build wioterminal

package main

import (
	"image/color"
	"machine"
	"strings"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/tinyfont/proggy"
	"tinygo.org/x/tinyterm"
)

var (
	display = ili9341.NewSPI(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)

	terminal = tinyterm.NewTerminal(display)

	black = color.RGBA{0, 0, 0, 255}
	white = color.RGBA{255, 255, 255, 255}
	red   = color.RGBA{255, 0, 0, 255}
	blue  = color.RGBA{0, 0, 255, 255}
	green = color.RGBA{0, 255, 0, 255}

	font = &proggy.TinySZ8pt7b
)

func setupDisplay() {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 40000000,
	})

	display.Configure(ili9341.Config{})
	display.SetRotation(ili9341.Rotation270)

	terminal.Configure(&tinyterm.Config{
		Font:              font,
		FontHeight:        10,
		FontOffset:        6,
		UseSoftwareScroll: true,
	})

	machine.LCD_BACKLIGHT.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LCD_BACKLIGHT.High()

	display.FillScreen(black)

	showSplash()
}

func showSplash() {
	for _, line := range strings.Split(strings.TrimSuffix(logo, "\n"), "\n") {
		terminal.Write([]byte("\n" + line))
	}
}

func terminalOutput(s string) {
	println(s)
	terminal.Write([]byte("\n" + s))
}
