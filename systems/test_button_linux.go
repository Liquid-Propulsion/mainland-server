/// TODO: This needs to be implemented, really wishing golang has rust macros here.
package systems

import (
	"log"

	"github.com/stianeikeland/go-rpio/v4"
)

type TestButton struct {
	hasRPIO       bool
	testButtonPin rpio.Pin
}

func NewTestButton() *TestButton {
	button := new(TestButton)
	err := rpio.Open()
	if err != nil {
		log.Printf("couldn't open test button, assuming this is in development mode: %s", err)
	} else {
		button.hasRPIO = true
		button.testButtonPin = rpio.Pin(10)
		button.testButtonPin.Input()
	}
	return button
}

func (engine *TestButton) ButtonHeld() bool {
	if engine.hasRPIO {
		return engine.testButtonPin.Read() == rpio.High
	}
	return true
}

func (engine *TestButton) HasRPIO() bool {
	return engine.hasRPIO
}
