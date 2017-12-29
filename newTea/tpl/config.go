package tpl

var GateConfigStr string = `package config

import (
	"time"

	"github.com/k4s/tea"
)

var (
	// gate conf
	Protocol          = "json" //"json"
	MaxConnNum        = 100
	WritingNum        = 2000
	MaxMsgLen  uint32 = 16384 //8192

	//tcp
	CTCPAddr     = "127.0.0.1:4444"
	STCPAddr     = "127.0.0.1:5555"
	LenMsgLen    = 2
	LittleEndian = true

	//websocket
	CWsAddr     = "127.0.0.1:4444"
	SWsAddr     = "127.0.0.1:5555"
	HTTPTimeout = 10 * time.Second
)

func GetOpts() tea.Options {
	var opts = make(tea.Options)
	opts.SetOption(tea.OptionConnNum, MaxConnNum)
	opts.SetOption(tea.OptionMaxRW, WritingNum)
	opts.SetOption(tea.OptionMaxMsgLen, MaxMsgLen)
	opts.SetOption(tea.OptioneMsgLen, LenMsgLen)
	opts.SetOption(tea.OptionLittleEndian, LittleEndian)
	return opts
}
`
var GameConfigStr string = `package config

import (
	"time"

	"github.com/k4s/tea"
)

var (
	// gate conf
	Protocol          = "json" //"json"
	MaxConnNum        = 100
	WritingNum        = 2000
	MaxMsgLen  uint32 = 16384 //8192

	//tcp
	GateTCPAddr  = "127.0.0.1:5555"
	LenMsgLen    = 2
	LittleEndian = false

	//websocket
	GateWSAddr  = "127.0.0.1:5555"
	HTTPTimeout = 10 * time.Second

	//log
	LogLevel = ""
	LogPath  = ""

	// other conf
	MysqlAddr = "root:@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true"
)

func GetOpts() tea.Options {
	var opts = make(tea.Options)
	opts.SetOption(tea.OptionConnNum, MaxConnNum)
	opts.SetOption(tea.OptionMaxRW, WritingNum)
	opts.SetOption(tea.OptionMaxMsgLen, MaxMsgLen)
	opts.SetOption(tea.OptioneMsgLen, LenMsgLen)
	opts.SetOption(tea.OptionLittleEndian, LittleEndian)
	return opts
}
`
