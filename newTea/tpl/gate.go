package tpl

var GateMainStr string = `package main

import (
	"<<DIR>>/config"

	"github.com/k4s/tea/gate"
)

func main() {
	gate := gate.NewGate(config.Appconfig.ClientAddr, config.Appconfig.GameServeraddr)
	gate.SetOpts(config.GetOpts())
	gate.Run()
}

`
