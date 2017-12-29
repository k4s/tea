package tpl

var RouterStr string = `package router

import (
	"<<DIR>>/handle"
	"<<DIR>>/msg"
	"<<DIR>>/protocol"
)

func InitRouter() {
	protocol.Processor.SetHandler(&msg.Hello{}, handle.InfoHandle)
}
`
