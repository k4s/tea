package tpl

var HandleStr string = `package handle

import (
	"fmt"
	ms "<<DIR>>/msg"
	"<<DIR>>/protocol"

	"github.com/k4s/tea/message"
	"github.com/k4s/tea/network"
)

func InfoHandle(msg *message.Message, agent network.Agent) {
	jsonMsg, err := protocol.Processor.Unmarshal(msg.Body)
	if err != nil {
		fmt.Println(err)
	}
	m := jsonMsg.(*ms.Hello)
	fmt.Println("game:", m)
	reMsg := ms.Hello{
		Name: "kkk",
	}
	agent.EndHandle(msg, reMsg)
}`
