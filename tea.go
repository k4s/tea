package tea

import (
	"os"
	"os/signal"

	"github.com/k4s/tea/gate"
	"github.com/k4s/tea/log"
)

type Tea struct {
	gate *gate.Gate
}

func NewTea(gate *gate.Gate) *Tea {
	tea := &Tea{
		gate: gate,
	}
	return tea
}

func (T *Tea) Run() {
	log.Release("Tea Game Server running...")
	T.gate.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	T.gate.Destroy()
	log.Release("Tea closing by (signal: %v)", sig)
}
