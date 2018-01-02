package network

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/k4s/tea/log"
	"github.com/k4s/tea/message"
)

type Agent interface {
	Run()
	WriteMsg(*message.Message)
	RecvMsg() *message.Message
	SetID(uint32)
	GetID() uint32
	UserData() interface{}
	SetUserData(data interface{})
	EndHandle(*message.Message, interface{})
	SetOnCloseFunc(func(Agent))
	IsClose() bool
	Close()
}

type agent struct {
	conn      Conn
	connID    uint32
	recvMsg   chan *message.Message
	withID    bool
	userData  interface{}
	closeFunc func(Agent)
	isClose   chan bool
}

func NewAgent(conn Conn, withID bool) *agent {
	return &agent{conn: conn,
		withID:  withID,
		isClose: make(chan bool),
		recvMsg: make(chan *message.Message, 100),
	}
}

func (a *agent) Run() {
	for {
		var msg *message.Message
		var err error
		if a.withID {
			msg, err = a.conn.RecvMsgWithID()
		} else {
			msg, err = a.conn.RecvMsg()
		}

		if err != nil {
			a.isClose <- true
			log.Debug("read message: %v", err)
			break
		}
		fmt.Println("agent:RUN:", string(msg.Body))
		a.recvMsg <- msg

	}
}

func (a *agent) SetOnCloseFunc(f func(Agent)) {
	a.closeFunc = f
}

func (a *agent) onClose() {
	if a.closeFunc != nil {
		a.closeFunc(a)
	}
}

func (a *agent) SetID(id uint32) {
	a.connID = id
}

func (a *agent) GetID() uint32 {
	return a.connID
}

func (a *agent) WriteMsg(msg *message.Message) {
	a.conn.SendMsg(msg)
}

//Endhandle free old msg,and write msg to gate
func (a *agent) EndHandle(oldMsg *message.Message, msg interface{}) {
	defer oldMsg.Free()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
	}
	sz := len(data)
	var newMsg *message.Message
	if oldMsg.ConnID != 0 {
		newMsg = message.NewMessage(int(+8))
		newMsg.ConnID = oldMsg.ConnID
		newMsg.Body = newMsg.Body[0:sz]
	} else {
		newMsg = message.NewMessage(int(+4))
		newMsg.Body = newMsg.Body[0:sz]
	}
	newMsg.Body = data
	a.WriteMsg(newMsg)
}

func (a *agent) RecvMsg() *message.Message {
	return <-a.recvMsg
}

//Close will exec onClosefunc and close net.conn
func (a *agent) Close() {
	a.onClose()
	a.conn.Close()
}

func (a *agent) IsClose() bool {
	return <-a.isClose
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
