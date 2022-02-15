package systems

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

type NodeSystem struct {
	nodes map[uint]bool
}

func NewNodeSystem() *NodeSystem {
	node := new(NodeSystem)
	go node.run()
	node.Reset()
	return node
}

func (system *NodeSystem) Reset() {
	system.nodes = make(map[uint]bool)
	var nodes []types.IslandNode
	res := sql.Database.Find(&nodes)
	if res.Error != nil {
		log.Printf("Couldn't query for nodes: %s", res.Error)
	}
	for _, node := range nodes {
		system.nodes[node.ID] = false
	}
	err := canbackend.CurrentCANBackend.SendPing()
	if err != nil {
		log.Printf("Ping Error: %s", err)
	}
}

func (system *NodeSystem) run() {
	for {
		pong := <-canbackend.CurrentCANBackend.PongChannel()
		if _, ok := system.nodes[uint(pong.NodeId)]; ok {
			system.nodes[uint(pong.NodeId)] = true
		}
	}
}

func (system *NodeSystem) NodeOnline(id uint) bool {
	online, ok := system.nodes[id]
	return online && ok
}
