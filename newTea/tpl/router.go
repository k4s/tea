package tpl

var RouterStr string = `package router

import (
	"<<DIR>>/game"
	"<<DIR>>/msg"
	"<<DIR>>/msg/process"
)

func init() {
	process.Processor.SetHandler(&msg.Login{}, game.MsgLogin)
}

`
