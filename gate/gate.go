package gate

import (
	"time"

	"github.com/k4s/tea/network/agent"
	. "github.com/k4s/tea/network/net"
	"github.com/k4s/tea/network/protocol"
)

type Gate struct {
	MaxConnNum int
	WritingNum int
	MaxMsgLen  uint32
	Processor  protocol.Processor

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	//close
	CloseSig chan bool
	IsClose  chan bool
}

func (gate *Gate) Run() {
	var wsServer *WSServer
	if gate.WSAddr != "" {
		wsServer = new(WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.WritingNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.NewAgent = func(conn *WSConn) agent.InAgent {
			a := agent.NewAgent(conn, gate.Processor)
			a.OnInit()
			return a
		}
	}

	var tcpServer *TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.WritingNum = gate.WritingNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *TCPConn) agent.InAgent {
			a := agent.NewAgent(conn, gate.Processor)
			a.OnInit()
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-gate.CloseSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
	gate.IsClose <- true
}

func (gate *Gate) Destroy() {
	gate.CloseSig <- true
	<-gate.IsClose
}
