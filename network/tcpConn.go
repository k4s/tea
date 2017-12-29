package network

import (
	"encoding/binary"
	"io"
	"net"
	"sync"

	. "github.com/k4s/tea"
	"github.com/k4s/tea/message"
)

type ConnSet map[net.Conn]struct{}

type Conn interface {
	SendMsg(*message.Message) error
	RecvMsg() (*message.Message, error)
	RecvMsgWithID() (*message.Message, error)
	Send([]byte) error
	Recv() ([]byte, error)
	Close() error
	IsOpen() bool
}

type TCPConn struct {
	sync.Mutex
	c        net.Conn
	writeMsg chan *message.Message
	isOpen   bool
	opts     Options
}

func newAgentTCPConn(c net.Conn, opts Options) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.c = c
	tcpConn.writeMsg = make(chan *message.Message, 100)
	tcpConn.opts = make(Options)
	// tcpConn.msgParser = msgParser
	for n, v := range opts {
		switch n {
		case OptionMinMsgLen,
			OptionMaxMsgLen,
			OptionMsgNum,
			OptionLittleEndian:
			tcpConn.opts.SetOption(n, v)
		}
	}
	if maxRW, err := tcpConn.opts.GetOption(OptionMaxRW); err != nil {
		tcpConn.opts.SetOption(OptionMaxRW, 4096)
	} else {
		switch maxRW := maxRW.(type) {
		case uint32:
			if maxRW <= 0 {
				tcpConn.opts.SetOption(OptionMaxRW, 4096)
			}
		}
	}

	go func() {
		for m := range tcpConn.writeMsg {
			if m == nil {
				break
			}
			err := tcpConn.doWriteMsgWithID(m)
			if err != nil {
				break
			}
		}

		tcpConn.Close()

	}()

	return tcpConn
}

func newSessionTCPConn(c net.Conn, opts Options) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.c = c
	tcpConn.writeMsg = make(chan *message.Message, 100)
	tcpConn.opts = make(Options)
	// tcpConn.msgParser = msgParser
	for n, v := range opts {
		switch n {
		case OptionMinMsgLen,
			OptionMaxMsgLen,
			OptionMsgNum,
			OptionLittleEndian:
			tcpConn.opts.SetOption(n, v)
		}
	}
	if maxRW, err := tcpConn.opts.GetOption(OptionMaxRW); err != nil {
		tcpConn.opts.SetOption(OptionMaxRW, 4096)
	} else {
		switch maxRW := maxRW.(type) {
		case uint32:
			if maxRW <= 0 {
				tcpConn.opts.SetOption(OptionMaxRW, 4096)
			}
		}
	}

	go func() {
		for m := range tcpConn.writeMsg {
			if m == nil {
				break
			}
			err := tcpConn.doWriteMsg(m)
			if err != nil {
				break
			}
		}

		tcpConn.Close()

	}()

	return tcpConn
}

func (c *TCPConn) RecvMsg() (*message.Message, error) {
	var sz uint32
	var err error
	var msg *message.Message

	isLittleEndian, err := c.opts.GetOption(OptionLittleEndian)

	if err == nil && isLittleEndian.(bool) {
		if err = binary.Read(c.c, binary.LittleEndian, &sz); err != nil {
			return nil, err
		}
	} else {
		if err = binary.Read(c.c, binary.BigEndian, &sz); err != nil {
			return nil, err
		}
	}
	maxRW, err := c.opts.GetOption(OptionMaxRW)
	if err == nil {
		if sz < 0 || (maxRW.(uint32) > 0 && sz > maxRW.(uint32)) {
			return nil, ErrTooLong
		}
	}
	msg = message.NewMessage(int(sz))
	msg.Body = msg.Body[0:sz]
	if _, err = io.ReadFull(c.c, msg.Body); err != nil {
		msg.Free()
		return nil, err
	}
	return msg, nil
}

func (c *TCPConn) RecvMsgWithID() (*message.Message, error) {
	var sz uint32
	var connID uint32
	var err error
	var msg *message.Message

	isLittleEndian, err := c.opts.GetOption(OptionLittleEndian)

	if err == nil && isLittleEndian.(bool) {
		if err = binary.Read(c.c, binary.LittleEndian, &sz); err != nil {
			return nil, err
		}
		if err = binary.Read(c.c, binary.LittleEndian, &connID); err != nil {
			return nil, err
		}
	} else {
		if err = binary.Read(c.c, binary.BigEndian, &sz); err != nil {
			return nil, err
		}
		if err = binary.Read(c.c, binary.BigEndian, &connID); err != nil {
			return nil, err
		}
	}
	maxRW, err := c.opts.GetOption(OptionMaxRW)
	if err == nil {
		if sz < 0 || (maxRW.(uint32) > 0 && sz > maxRW.(uint32)) {
			return nil, ErrTooLong
		}
	}
	msg = message.NewMessage(int(sz))
	msg.ConnID = connID
	msg.Body = msg.Body[0 : sz-4]
	if _, err = io.ReadFull(c.c, msg.Body); err != nil {
		msg.Free()
		return nil, err
	}
	return msg, nil
}

func (c *TCPConn) doWriteMsg(msg *message.Message) error {
	sz := uint32(len(msg.Body))
	m := make([]byte, sz+4)

	isLittleEndian, err := c.opts.GetOption(OptionLittleEndian)
	if err == nil && isLittleEndian.(bool) {
		binary.LittleEndian.PutUint32(m, sz)
	} else {
		binary.BigEndian.PutUint32(m, sz)
	}
	copy(m[4:], msg.Body)
	if _, err := c.c.Write(m); err != nil {
		return err
	}
	msg.Free()
	return nil
}

func (c *TCPConn) doWriteMsgWithID(msg *message.Message) error {
	sz := uint32(len(msg.Body) + 4)
	m := make([]byte, sz+4)
	isLittleEndian, err := c.opts.GetOption(OptionLittleEndian)
	if err == nil && isLittleEndian.(bool) {
		binary.LittleEndian.PutUint32(m, sz)
		binary.LittleEndian.PutUint32(m[4:], msg.ConnID)
	} else {
		binary.BigEndian.PutUint32(m, sz)
		binary.BigEndian.PutUint32(m[4:], msg.ConnID)
	}
	copy(m[8:], msg.Body)
	if _, err := c.c.Write(m); err != nil {
		return err
	}
	msg.Free()
	return nil
}

func (c *TCPConn) SendMsg(msg *message.Message) error {
	c.writeMsg <- msg
	return nil
}

func (c *TCPConn) Send(b []byte) error {
	msg := message.NewMessage(len(b))
	msg.Body = append(msg.Body, b...)
	return c.SendMsg(msg)
}

func (c *TCPConn) Recv() ([]byte, error) {
	msg, err := c.RecvMsg()
	if err != nil {
		return nil, err
	}
	b := make([]byte, 0, len(msg.Body))
	b = append(b, msg.Body...)
	msg.Free()
	return b, nil
}

func (c *TCPConn) IsOpen() bool {
	return c.isOpen
}

func (c *TCPConn) Close() error {
	c.Lock()
	if c.isOpen {
		c.isOpen = false
		for {
			if len(c.writeMsg) == 0 {
				break
			} else {
				m := <-c.writeMsg
				m.Free()
			}
		}
		close(c.writeMsg)
	}
	c.Unlock()
	return c.c.Close()
}
