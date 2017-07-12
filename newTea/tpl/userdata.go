package tpl

var UserdataStr string = `package msg

//用户登录的基本信息
//这些信息将保存在gate.Agent.UserData
//信息包含用户名，是否认证登录
type UserData struct {
	Uid      string
	Account  string
	Password string
	Auth     bool
}

`
