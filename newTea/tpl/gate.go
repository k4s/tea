package tpl

var GateStr string = `package gate

import (
	"github.com/k4s/tea/gate"
	"<<DIR>>/msg/process"
	"<<DIR>>/config"
)

var Gate = &gate.Gate{
	MaxConnNum:   config.MaxConnNum,
	WritingNum:   config.WritingNum,
	MaxMsgLen:    config.MaxMsgLen,
	WSAddr:       config.WSAddr,
	HTTPTimeout:  config.HTTPTimeout,
	TCPAddr:      config.TCPAddr,
	LenMsgLen:    config.LenMsgLen,
	LittleEndian: config.LittleEndian,
	Processor:    process.Processor,
}
`
