package tpl

var MainStr string = `package main

import (
	"github.com/k4s/tea"
	"<<DIR>>/gate"
	_ "<<DIR>>/router"
)

func main() {
	tea := tea.NewTea(gate.Gate)
	tea.Run()
}
`
