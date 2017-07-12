package net

import (
	"net"
	"sync"
	"time"

	"github.com/k4s/tea/log"
	"github.com/k4s/tea/network/agent"
)

type TCPServer struct {
	//server address："127.0.0.1:8080"
	Addr string
	//Max connections number
	MaxConnNum int
	//writing byte number
	WritingNum int
	//connections number
	conns ConnSet
	//connections mutex
	mutexConns sync.Mutex
	//listen waitGroup
	wglisten sync.WaitGroup
	//writing waitGroup
	wgConns  sync.WaitGroup
	Timeout  time.Duration
	listener net.Listener
	NewAgent func(*TCPConn) agent.InAgent

	// msg parser
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
	msgParser    *MsgParser
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
	if server.MaxConnNum <= 0 {
		server.MaxConnNum = 100
		log.Release("invalid MaxConnNum, reset to %v", server.MaxConnNum)
	}
	if server.WritingNum <= 0 {
		server.WritingNum = 100
		log.Release("invalid WritingNum, reset to %v", server.WritingNum)
	}
	if server.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}
	server.listener = listener
	server.conns = make(ConnSet)

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(server.LenMsgLen, server.MinMsgLen, server.MaxMsgLen)
	msgParser.SetByteOrder(server.LittleEndian)
	server.msgParser = msgParser

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
		if len(server.conns) >= server.MaxConnNum {
			server.mutexConns.Unlock()
			conn.Close()
			log.Debug("too many connections")
			continue
		}

		//Collect all the connections
		server.conns[conn] = struct{}{}
		server.mutexConns.Unlock()

		server.wgConns.Add(1)

		tcpConn := newTCPConn(conn, server.WritingNum, server.msgParser)
		agent := server.NewAgent(tcpConn)

		go func() {
			agent.Run()
			tcpConn.Close()

			server.mutexConns.Lock()
			delete(server.conns, conn)
			server.mutexConns.Unlock()
			agent.OnClose()

			server.wgConns.Done()
		}()
	}

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
