package network

import (
	"net"
	"sync"
	"time"

	. "github.com/k4s/tea"
	"github.com/k4s/tea/log"
)

type TCPClient struct {
	sync.Mutex
	Addr string

	conns ConnSet
	//connections mutex
	mutexConns sync.Mutex
	//writing waitGroup
	wgConns sync.WaitGroup

	NewAgent func(*TCPConn, bool) Agent
	agent    Agent

	isClose bool
	withID  bool
	opts    Options

	// msgParser *MsgParser
}

func (client *TCPClient) Start() {
	client.init()

	go client.connect()
}

func (client *TCPClient) init() {
	client.Lock()
	defer client.Unlock()
	client.opts = make(Options)
	if connNum, err := client.opts.GetOption(OptionConnNum); err != nil {
		switch connNum := connNum.(type) {
		case int:
			if connNum <= 0 {
				client.opts.SetOption(OptionConnNum, 1)
			}
		}
	} else {
		client.opts.SetOption(OptionConnNum, 1)
	}

	if connInterval, err := client.opts.GetOption(OptionConnInterval); err != nil {
		client.opts.SetOption(OptionConnInterval, 3*time.Second)
	} else {
		switch connInterval := connInterval.(type) {
		case int:
			if connInterval <= 0 {
				client.opts.SetOption(OptionConnInterval, 3*time.Second)
			}
		}
	}

	if msgNum, err := client.opts.GetOption(OptionMsgNum); err != nil {
		client.opts.SetOption(OptionMsgNum, 100)
	} else {
		switch msgNum := msgNum.(type) {
		case int:
			if msgNum <= 0 {
				client.opts.SetOption(OptionMsgNum, 100)
			}
		}
	}
	if isLittleEndian, err := client.opts.GetOption(OptionLittleEndian); err != nil {
		client.opts.SetOption(OptionLittleEndian, true)
	} else {
		switch isLittleEndian := isLittleEndian.(type) {
		case bool:
			if !isLittleEndian {
				client.opts.SetOption(OptionLittleEndian, true)
			}
		}
	}

	if client.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}
	if client.conns != nil {
		log.Fatal("client is running")
	}

	client.conns = make(ConnSet)
	client.isClose = false

}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil || client.isClose {
			return conn
		}
		log.Release("connect to %v error: %v", client.Addr, err)

		connInterval, ok := client.opts.GetOption(OptionConnInterval)
		if ok != nil {
			connInterval = 3 * time.Second
		}
		time.Sleep(connInterval.(time.Duration))
		continue
	}
}

func (client *TCPClient) connect() {

	conn := client.dial()
	if conn == nil {
		return
	}

	client.Lock()
	if client.isClose {
		client.Unlock()
		conn.Close()
		return
	}
	client.conns[conn] = struct{}{}
	client.Unlock()

	var tcpConn *TCPConn
	if client.withID {
		tcpConn = newAgentTCPConn(conn, client.opts)
	} else {
		tcpConn = newSessionTCPConn(conn, client.opts)
	}
	agent := client.NewAgent(tcpConn, client.withID)
	client.agent = agent
	agent.Run()

	// cleanup
	tcpConn.Close()
	client.Lock()
	delete(client.conns, conn)
	client.Unlock()
	// agent.OnClose()
	agent.Close()
}

func (client *TCPClient) SetWithID(isbool bool) {
	client.withID = isbool
}

func (client *TCPClient) SetOpts(opts Options) {
	client.opts = opts
}

func (client *TCPClient) GetAgent() Agent {
	return client.agent
}

func (client *TCPClient) Close() {
	client.Lock()
	client.isClose = true
	for conn := range client.conns {
		conn.Close()
	}
	client.conns = nil
	client.Unlock()

}
