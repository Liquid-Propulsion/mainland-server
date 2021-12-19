package systems

type TestButton struct {
}

func NewTestButton() *TestButton {
	button := new(TestButton)
	return button
}

func (engine *TestButton) ButtonHeld() bool {
	return true
}

func (engine *TestButton) HasRPIO() bool {
	return false
}
