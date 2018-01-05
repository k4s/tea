package tpl

var GateConfigJson string = `{
    "Protocal": "json",
    "MaxConnNum": 100,
    "MaxRW": 2000,
    "MaxMsgLen": 8192,
    "ClientAddr": "127.0.0.1:4444",
    "GameServeraddr": "127.0.0.1:5555",
    "HTTPTimeout": 10,
    "LittleEndian": true
}
`

var GateConfigStr string = `package config

import (
	"path/filepath"

	"github.com/k4s/tea"
)

func init() {
	configFilePath, _ := filepath.Abs("./config/config.json")
	err := tea.InitConfigFile(configFilePath, &Appconfig)
	if err != nil {
		panic(err)
	}
}

//Appconfig 应用配置
var Appconfig struct {
	Protocol       string ` + "`" + `json:"Protocal"` + "`" + `
	MaxRW          int    ` + "`" + `json:"MaxRW"` + "`" + `
	MaxConnNum     int    ` + "`" + `json:"MaxConnNum"` + "`" + `
	MaxMsgLen      int    ` + "`" + `json:"MaxMsgLen"` + "`" + `
	ClientAddr     string ` + "`" + `json:"ClientAddr"` + "`" + `
	GameServeraddr string ` + "`" + `json:"GameServeraddr"` + "`" + `
	HTTPTimeout    int    ` + "`" + `json:"HTTPTimeout"` + "`" + `
	LittleEndian   bool   ` + "`" + `json:"LittleEndian"` + "`" + `
}

func GetOpts() tea.Options {
	var opts = make(tea.Options)
	opts.SetOption(tea.OptionConnNum, Appconfig.MaxConnNum)
	opts.SetOption(tea.OptionMaxRW, Appconfig.MaxRW)
	opts.SetOption(tea.OptionMaxMsgLen, Appconfig.MaxMsgLen)
	opts.SetOption(tea.OptionLittleEndian, Appconfig.LittleEndian)
	return opts
}

`
var GameConfigJson string = `{
    "Protocal": "json",
    "MaxConnNum": 100,
    "MaxRW": 2000,
    "MaxMsgLen": 8192,
    "GateTCPAddr": [
        "127.0.0.1:5555"
    ],
    "GateWSAddr": [
        "127.0.0.1:5555"
    ],
    "HTTPTimeout": 10,
    "LittleEndian": true,
    "MysqlAddr": "root:@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true"
}
`

var GameConfigStr string = `package config

import (
	"path/filepath"

	"github.com/k4s/tea"
)

func init() {
	configFilePath, _ := filepath.Abs("./config/config.json")
	err := tea.InitConfigFile(configFilePath, &Appconfig)
	if err != nil {
		panic(err)
	}
}

//Appconfig 应用配置
var Appconfig struct {
	Protocol     string   ` + "`" + `json:"Protocal"` + "`" + `
	MaxRW        int      ` + "`" + `json:"MaxRW"` + "`" + `
	MaxConnNum   int      ` + "`" + `json:"MaxConnNum"` + "`" + `
	MaxMsgLen    int      ` + "`" + `json:"MaxMsgLen"` + "`" + `
	GateTCPAddr  []string ` + "`" + `json:"GateTCPAddr"` + "`" + `
	GateWSAddr   []string ` + "`" + `json:"GateWSAddr"` + "`" + `
	HTTPTimeout  int      ` + "`" + `json:"HTTPTimeout"` + "`" + `
	LittleEndian bool     ` + "`" + `json:"LittleEndian"` + "`" + `
	MysqlAddr    string   ` + "`" + `json:"MysqlAddr"` + "`" + `
}

func GetOpts() tea.Options {
	var opts = make(tea.Options)
	opts.SetOption(tea.OptionConnNum, Appconfig.MaxConnNum)
	opts.SetOption(tea.OptionMaxRW, Appconfig.MaxRW)
	opts.SetOption(tea.OptionMaxMsgLen, Appconfig.MaxMsgLen)
	opts.SetOption(tea.OptionLittleEndian, Appconfig.LittleEndian)
	return opts
}
`
