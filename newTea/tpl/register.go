package tpl

var RegisterStr string = `package Register

import (
	"<<DIR>>/msg"
	"<<DIR>>/protocol"
)

func InitRegister() {
	protocol.Processor.Register(&msg.Hello{})
}`
