package tpl

var ProcessStr string = `package process

import (
	"log"

	"github.com/k4s/tea/network/protocol"
	"<<DIR>>/config"
)

var Processor protocol.Processor

func init() {
	switch config.Protocol {
	case "json":
		Processor = protocol.NewJson()
	case "protobuf":
		Processor = protocol.NewProto()
	default:
		log.Fatal("unknown Protocol: %v", config.Protocol)
	}

}

`
