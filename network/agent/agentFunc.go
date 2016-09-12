package agent

import (
	"github.com/k4s/tea/network"
)

var (
	CloseFunc func(network.ExAgent)
	InitFunc  func(network.ExAgent)
)
