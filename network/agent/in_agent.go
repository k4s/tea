package agent

type InAgent interface {
	Run()
	OnInit()
	OnClose()
}
