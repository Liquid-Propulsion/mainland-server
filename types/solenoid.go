package types

type Solenoid struct {
	Model
	Name        string
	Description string
}

func (node *Solenoid) IsNode() {}
