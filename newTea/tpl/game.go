package tpl

var GameStr string = `package game

import (
	"fmt"

	"github.com/k4s/tea/network"
	ms "<<DIR>>/msg"
)

func MsgHello(msg interface{}, agent network.ExAgent) {
	m := msg.(*ms.Hello)
	fmt.Println("Hello,", m.Name)
	hi := &ms.Hello{Name: "kas"}
	agent.WriteMsg(hi)
}

func MsgLogin(msg interface{}, agent network.ExAgent) {
	m := msg.(*ms.Login)
	if m.User == "kas" && m.Password == "123456" {
		userdata := ms.UserData{}

		userdata.Account = m.User
		userdata.Password = m.Password
		userdata.Auth = true
		agent.SetUserData(userdata)

		hi := &ms.Hello{Name: "kas"}
		agent.WriteMsg(hi)

		user := agent.UserData().(ms.UserData)
		fmt.Println("user is auth :", user)

		return
	}
}`
