package gameserver

import (
	"fmt"
	"os"
	"os/signal"

	. "github.com/k4s/tea"
	"github.com/k4s/tea/log"
	"github.com/k4s/tea/network"
	"github.com/k4s/tea/protocol"
)

type Gameserver struct {
	TCPAddr   string
	opts      Options
	CloseChan chan bool
	IsClose   chan bool
	Processor protocol.Processor
}

func NewGameserver(addr string, processor protocol.Processor) *Gameserver {
	gate := &Gameserver{
		TCPAddr:   addr,
		opts:      make(Options),
		CloseChan: make(chan bool),
		IsClose:   make(chan bool),
		Processor: processor,
	}
	return gate
}

func (g *Gameserver) Run() {
	log.Release("Gameserver Server running by %s", g.TCPAddr)
	g.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	g.Stop()
	log.Release("Gameserver closing by (signal: %v)", sig)
}

func (g *Gameserver) Start() {
	var tcpClient *network.TCPClient
	if g.TCPAddr != "" {
		tcpClient = new(network.TCPClient)
		tcpClient.Addr = g.TCPAddr
		tcpClient.NewAgent = func(conn *network.TCPConn, withID bool) network.Agent {
			a := network.NewAgent(conn, withID)
			return a
		}
		tcpClient.SetOpts(g.opts)
		tcpClient.SetWithID(true)
	}
	if tcpClient != nil {
		tcpClient.Start()
		go g.WaitAgent(tcpClient)

	}
	<-g.CloseChan
	if tcpClient != nil {
		tcpClient.Close()
	}
	g.IsClose <- true
}

func (g *Gameserver) Stop() {
	g.CloseChan <- true
	<-g.IsClose
}

func (g *Gameserver) WaitAgent(client *network.TCPClient) {
	var agent network.Agent
	for {
		agent = client.GetAgent()
		if agent != nil {
			fmt.Println("取出不为空")
			break
		}
		fmt.Println("取出为空")
	}
	go g.RunAgent(agent)
}

func (g *Gameserver) RunAgent(agent network.Agent) {
	for {
		msg := agent.RecvMsg()
		g.Processor.Route(msg, agent)

	}
}
