package gate

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"

	. "github.com/k4s/tea"
	"github.com/k4s/tea/log"
	"github.com/k4s/tea/network"
)

type Gate struct {
	clientAddr string
	serverAddr string
	cNum       uint32
	cAgents    map[uint32]network.Agent
	sAgent     map[network.Agent]struct{}
	sync.Mutex
	opts      Options
	CloseChan chan bool
	IsClose   chan bool
}

func NewGate(cAddr, sAddr string) *Gate {
	gate := &Gate{
		clientAddr: cAddr,
		serverAddr: sAddr,
		cAgents:    make(map[uint32]network.Agent),
		cNum:       0,
		sAgent:     make(map[network.Agent]struct{}),
		opts:       make(Options),
		CloseChan:  make(chan bool),
		IsClose:    make(chan bool),
	}
	return gate
}

func (g *Gate) SetOpts(opts Options) {
	g.opts = opts
}

func (g *Gate) Run() {
	log.Release("Gateway External Server Running : %s", g.clientAddr)
	log.Release("Gateway Internal Server Running : %s", g.serverAddr)
	g.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	g.Stop()
	log.Release("Gateway closing by (signal: %v)", sig)
}

func (g *Gate) Start() {
	var CServer *network.TCPServer
	var SServer *network.TCPServer
	if g.clientAddr != "" {
		CServer = new(network.TCPServer)
		CServer.Addr = g.clientAddr
		CServer.NewAgent = func(conn *network.TCPConn, withID bool) network.Agent {
			a := network.NewAgent(conn, withID)
			return a
		}
		CServer.SetOpts(g.opts)

	}
	if g.serverAddr != "" {
		SServer = new(network.TCPServer)
		SServer.Addr = g.serverAddr
		SServer.NewAgent = func(conn *network.TCPConn, withID bool) network.Agent {
			a := network.NewAgent(conn, withID)
			return a
		}
		CServer.SetOpts(g.opts)
		SServer.SetWithID(true)
	}
	if CServer != nil {
		CServer.Start()

		go g.waitSession(CServer)
	}
	if SServer != nil {
		SServer.Start()

		go g.waitAgent(SServer)
	}
	<-g.CloseChan
	if CServer != nil {
		CServer.Close()
	}
	if SServer != nil {
		SServer.Close()
	}
	g.IsClose <- true
}
func (g *Gate) Stop() {
	g.CloseChan <- true
	<-g.IsClose
}

func (g *Gate) waitSession(server *network.TCPServer) {
	for {
		agent := server.GetAgent()
		for {
			atomic.AddUint32(&g.cNum, 1)
			if _, ok := g.cAgents[g.cNum]; !ok {
				g.Lock()
				g.cAgents[g.cNum] = agent
				g.Unlock()
				agent.SetID(g.cNum)
				go g.stopSession(agent)
				break
			}
		}

		go g.runSession(agent)
	}
}

func (g *Gate) stopSession(agent network.Agent) {
	agent.IsClose()
	g.Lock()
	delete(g.cAgents, agent.GetID())
	g.Unlock()
}

func (g *Gate) runSession(agent network.Agent) {
	for {
		msg := agent.RecvMsg()
		msg.ConnID = agent.GetID()

		//gate To gmaeserver
		for a := range g.sAgent {
			a.WriteMsg(msg)
			break
		}
		msg.Free()

	}
}

func (g *Gate) waitAgent(server *network.TCPServer) {
	for {
		agent := server.GetAgent()
		g.Lock()
		g.sAgent[agent] = struct{}{}
		g.Unlock()
		go g.stopAgent(agent)
		go g.runAgent(agent)
	}
}
func (g *Gate) stopAgent(agent network.Agent) {
	agent.IsClose()
	g.Lock()
	delete(g.sAgent, agent)
	g.Unlock()
}

func (g *Gate) runAgent(agent network.Agent) {
	for {
		msg := agent.RecvMsg()
		if a, ok := g.cAgents[msg.ConnID]; ok {
			a.WriteMsg(msg)
			msg.Free()
		}

	}
}
