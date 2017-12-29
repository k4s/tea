package protocol

import (
	"github.com/k4s/tea/message"
	"github.com/k4s/tea/network"
)

type Processor interface {
	// must goroutine safe
	Route(msg *message.Message, agent network.Agent) error
	//must goroutine safe
	SetHandler(msg interface{}, msgHandler MsgHandler)
	//must goroutine safe
	Register(msg interface{})
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error)
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
}
