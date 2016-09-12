package protocol

import (
	"reflect"

	"github.com/k4s/tea/network"
)

type MsgInfo struct {
	msgType    reflect.Type
	msgHandler MsgHandler
}

type MsgHandler func(interface{}, network.ExAgent)
