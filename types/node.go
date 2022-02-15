package types

type IslandNode struct {
	Model
	Name        string
	Description string
}

func (node *IslandNode) IsNode() {}
