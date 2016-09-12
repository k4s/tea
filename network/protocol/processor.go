package protocol

import (
	"github.com/k4s/tea/network"
)

type Processor interface {
	// must goroutine safe
	Route(msg interface{}, agent network.ExAgent) error
	//must goroutine safe
	SetHandler(msg interface{}, msgHandler MsgHandler)
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error)
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error)
}
