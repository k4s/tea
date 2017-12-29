package tpl

var GameMainStr string = `package main

import (
	"<<DIR>>/config"
	"<<DIR>>/protocol"
	. "<<DIR>>/register"
	. "<<DIR>>/router"

	"github.com/k4s/tea/gameserver"
)

func init() {
	InitRegister()
	InitRouter()
}

func main() {
	game := gameserver.NewGameserver(config.GateTCPAddr, protocol.Processor)
	game.Run()
}`
