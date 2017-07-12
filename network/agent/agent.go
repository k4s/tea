package agent

import (
	"reflect"

	"github.com/k4s/tea/log"
	"github.com/k4s/tea/network"
	"github.com/k4s/tea/network/protocol"
)

type agent struct {
	conn      network.Conn
	processor protocol.Processor
	userData  interface{}
	closeFunc func(network.ExAgent)
	initFunc  func(network.ExAgent)
}

func NewAgent(conn network.Conn, processor protocol.Processor) *agent {
	return &agent{conn: conn, processor: processor, closeFunc: CloseFunc, initFunc: InitFunc}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.processor != nil {
			msg, err := a.processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnInit() {
	if a.initFunc != nil {
		a.initFunc(a)
	}
}

func (a *agent) OnClose() {
	if a.closeFunc != nil {
		a.closeFunc(a)
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.processor != nil {
		data, err := a.processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		a.conn.WriteMsg(data...)
	}
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
