package tpl

var ProtocolStr string = `package protocol

import (
	"log"

	"<<DIR>>/config"

	"github.com/k4s/tea/protocol"
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
