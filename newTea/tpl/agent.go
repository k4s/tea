package tpl

var AgentStr string = `package game

import (
	"fmt"

	"github.com/k4s/tea/network"
	"github.com/k4s/tea/network/agent"
)

func init() {
	agent.CloseFunc = closeAgent

	agent.InitFunc = initAgent
}

func closeAgent(agent network.ExAgent) {
	fmt.Println("CloseFunc")
}

func initAgent(agent network.ExAgent) {
	fmt.Println("InitFunc")
}`
