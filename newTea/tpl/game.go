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
	games := make([]*gameserver.Gameserver, 0)
	for _, addr := range config.Appconfig.GateTCPAddr {
		game := gameserver.NewGameserver(addr, protocol.Processor)
		games = append(games, game)
	}
	gameserver.GameRun(games)
}`
