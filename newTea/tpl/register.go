package tpl

var RegisterStr string = `package msg

import (
	"<<DIR>>/msg/process"
)

func init() {
	process.Processor.Register(&Login{})
	process.Processor.Register(&Hello{})
}

`
