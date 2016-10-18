package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/k4s/tea/network"
	"github.com/name5566/leaf/log"
)

// --------------------
// len | json message |
// --------------------
type JsonProcessor struct {
	msgInfo map[string]*MsgInfo
}

func NewJson() *JsonProcessor {
	p := new(JsonProcessor)
	p.msgInfo = make(map[string]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *JsonProcessor) Register(msg interface{}) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := msgType.Elem().Name()

	if msgID == "" {
		log.Fatal("unnamed json message")
	}
	if _, ok := p.msgInfo[msgID]; ok {
		log.Fatal("message %v is already registered", msgID)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[msgID] = i
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *JsonProcessor) SetHandler(msg interface{}, msgHandler MsgHandler) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}

	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgHandler = msgHandler
}

// goroutine safe
func (p *JsonProcessor) Route(msg interface{}, agent network.ExAgent) error {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		return fmt.Errorf("message %v not registered", msgID)
	}

	if i.msgHandler != nil {
		i.msgHandler(msg, agent)
	}
	return nil
}

// goroutine safe
func (p *JsonProcessor) Unmarshal(data []byte) (interface{}, error) {
	fmt.Println(string(data))
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	if len(m) != 1 {
		return nil, errors.New("invalid json data")
	}

	for msgID, data := range m {
		i, ok := p.msgInfo[msgID]
		if !ok {
			return nil, fmt.Errorf("message %v not registered", msgID)
		}

		// msg
		msg := reflect.New(i.msgType.Elem()).Interface()
		return msg, json.Unmarshal(data, msg)
	}

	panic("bug")
}

// goroutine safe
func (p *JsonProcessor) Marshal(msg interface{}) ([][]byte, error) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return nil, errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	if _, ok := p.msgInfo[msgID]; !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}

	// data
	m := map[string]interface{}{msgID: msg}
	data, err := json.Marshal(m)
	return [][]byte{data}, err
}
