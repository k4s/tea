package tpl

var GateMainStr string = `package main

import (
	"<<DIR>>/config"

	"github.com/k4s/tea/gate"
)

func main() {
	gate := gate.NewGate(config.CTCPAddr, config.STCPAddr)
	gate.SetOpts(config.GetOpts())
	gate.Run()
}
`
