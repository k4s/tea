package network

import (
	"net"
	"sync"
	"time"

	. "github.com/k4s/tea"
	"github.com/k4s/tea/log"
)

type TCPServer struct {
	//server address："127.0.0.1:8080"
	Addr string

	conns ConnSet
	//connections mutex
	mutexConns sync.Mutex
	//listen waitGroup
	wglisten sync.WaitGroup
	//writing waitGroup
	wgConns sync.WaitGroup

	listener net.Listener
	NewAgent func(*TCPConn, bool) Agent
	Agents   chan Agent

	withID bool
	opts   Options
}

func (server *TCPServer) Start() {
	server.init()
	go server.run()
}

func (server *TCPServer) init() {
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal("%v", err)
	}
	server.Agents = make(chan Agent, 100)
	server.opts = make(Options)
	for n, v := range server.opts {
		switch n {
		case OptionMinMsgLen,
			OptionMaxMsgLen,
			OptionConnNum,
			OptionMsgNum,
			OptionLittleEndian:
			server.opts.SetOption(n, v)
		}
	}

	if _, err := server.opts.GetOption(OptionConnNum); err != nil {
		server.opts.SetOption(OptionConnNum, 100)
	}

	if _, err := server.opts.GetOption(OptionMsgNum); err != nil {
		server.opts.SetOption(OptionMsgNum, 100)
	}

	if _, err := server.opts.GetOption(OptionLittleEndian); err != nil {
		server.opts.SetOption(OptionLittleEndian, true)
	}

	if server.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}
	server.listener = listener
	server.conns = make(ConnSet)

}

func (server *TCPServer) run() {
	server.wglisten.Add(1)
	defer server.wglisten.Done()

	var tempDelay time.Duration
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Release("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		server.mutexConns.Lock()
		connNum, _ := server.opts.GetOption(OptionConnNum)
		if len(server.conns) >= connNum.(int) {
			server.mutexConns.Unlock()
			conn.Close()
			log.Debug("too many connections")
			continue
		}

		//Collect all the connections
		server.conns[conn] = struct{}{}
		server.mutexConns.Unlock()

		server.wgConns.Add(1)

		var tcpConn *TCPConn
		if server.withID {
			tcpConn = newAgentTCPConn(conn, server.opts)
		} else {
			tcpConn = newSessionTCPConn(conn, server.opts)
		}

		agent := server.NewAgent(tcpConn, server.withID)

		go func() {
			server.Agents <- agent
			agent.Run()
			agent.Close()

			server.mutexConns.Lock()
			delete(server.conns, conn)
			server.mutexConns.Unlock()

			server.wgConns.Done()
		}()
	}

}

func (server *TCPServer) SetWithID(isbool bool) {
	server.withID = isbool
}

func (server *TCPServer) SetOpts(opts Options) {
	server.opts = opts
}

func (server *TCPServer) GetAgent() Agent {
	return <-server.Agents
}

func (server *TCPServer) Close() {
	server.listener.Close()
	server.wglisten.Wait()

	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()
	server.wgConns.Wait()
}
