package network

type ExAgent interface {
	WriteMsg(msg interface{})
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	OnInit()
	OnClose()
}
