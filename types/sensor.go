package types

type Sensor struct {
	Model
	Name          string
	Description   string
	TransformCode string
}

func (node *Sensor) IsNode() {}
