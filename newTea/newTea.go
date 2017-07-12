package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/k4s/tea/log"
	"github.com/k4s/tea/newTea/tpl"
)

var appname string
var crupath, _ = os.Getwd()

func main() {
	if len(os.Args) < 2 {
		fmt.Println("newTea new appname")
		return
	}
	switch os.Args[1] {
	case "new":
		initVar(os.Args[2])
		appname = os.Args[2]
		newapp(appname)
	}
}

func newapp(args string) {
	log.Debug("[tea] Create a tea project named `%s` in the `%s` path.", appname, crupath)
	if isExist(crupath) {
		log.Debug("[tea] The project path has conflic, do you want to build in: %s\n", crupath)
		log.Debug("[tea] Do you want to overwrite it? [(yes|no) or (y|n)]  ")
		if !askForConfirmation() {
			log.Fatal("[tea] Cancel...")
			return
		}
	}

	log.Debug("[tea] Start create project...")

	//生成配置文件
	makedir("config")
	writefile(crupath+"/config/config.go", replaceAppname(tpl.ConfigStr))

	//生成游戏逻辑文件
	makedir("game")
	writefile(crupath+"/game/game.go", replaceAppname(tpl.GameStr))
	writefile(crupath+"/game/agent.go", replaceAppname(tpl.AgentStr))

	//生成网关文件
	makedir("gate")
	writefile(crupath+"/gate/gate.go", replaceAppname(tpl.GateStr))

	//生成日志文件
	makedir("log/logdata")
	writefile(crupath+"/log/log.go", replaceAppname(tpl.LogStr))

	//生成消息文件
	makedir("msg/process")
	writefile(crupath+"/msg/process/process.go", replaceAppname(tpl.ProcessStr))
	writefile(crupath+"/msg/msg.go", replaceAppname(tpl.MsgStr))
	writefile(crupath+"/msg/register.go", replaceAppname(tpl.RegisterStr))
	writefile(crupath+"/msg/userdata.go", replaceAppname(tpl.UserdataStr))

	//生成路由文件
	makedir("router")
	writefile(crupath+"/router/router.go", replaceAppname(tpl.RouterStr))

	//生成main文件
	writefile(crupath+"/main.go", replaceAppname(tpl.MainStr))

	log.Debug("[tea] Create was successful")

	if err := os.Chdir(crupath); err != nil {
		log.Fatal("[tea] Create project fail: %v", err)
	}

}
func initVar(args string) {
	var dir string
	dir, appname = filepath.Split(args)
	if dir != "" {
		crupath = filepath.Join(dir, appname)
	} else {
		crupath = filepath.Join(crupath, appname)
	}

	var err error
	crupath = strings.TrimSpace(crupath)
	crupath, err = filepath.Abs(crupath)
	if err != nil {
		log.Fatal("[tea] Create project fail: %s", err)
	}
	crupath = strings.Replace(crupath, `\`, `/`, -1)
	crupath = strings.TrimRight(crupath, "/") + "/"
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

/**
 * 生成目录
 */
func makedir(dirname string) {
	err := os.MkdirAll(crupath+dirname, os.ModePerm)
	if err != nil {
		log.Fatal("[tea] create has a error :%s", err)
		return
	}
}

/**
 * 写文件
 */
func writefile(filename string, writeStr string) {
	var f *os.File
	var err error

	if checkFileIsExist(filename) { //如果文件存在
		f, err = os.OpenFile(filename, os.O_RDWR, 0666) //打开文件
		if err != nil {
			log.Fatal("[tea] create has a error :%s", err)
			return
		}
	} else {
		f, err = os.Create(filename) //创建文件
		if err != nil {
			log.Fatal("[tea] create has a error :%s", err)
			return
		}
	}
	_, err = io.WriteString(f, writeStr) //写入文件(字符串)
	if err != nil {
		log.Fatal("[tea] create has a error :%s", err)
		return
	}
}

/**
 *判断当前路径是否存在
 */
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

/**
 *替换appname
 */
func replaceAppname(inStr string) string {
	return strings.Replace(inStr, "<<DIR>>", appname, -1)
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err.Error())
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}
