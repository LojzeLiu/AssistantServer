package main

import (
	"AssistantServer/DisplyPck"
	"Common"
	"flag"
	"fmt"
	"net/http"
	"strconv"
)

var ConfFile = flag.String("f", "./WevServer.conf", "configer file")

func main() {
	//解析配置
	conf := &Common.Configer{}
	if err := conf.Init(ConfFile); err != nil {
		fmt.Println("Configer initialization Error:", err)
		return
	}
	defer conf.Destroy()

	//获取日志配置
	logConfs, err := conf.GetConf("LOG")
	if err != nil {
		fmt.Println("Get log configuer Error:", err)
		return
	}
	DebugLevel, err := strconv.Atoi(logConfs["Level"])
	if err != nil {
		fmt.Println("Log level Error:", err)
		return
	}

	//初始化日志
	if err := Common.SetLogger(logConfs["Path"], logConfs["AppName"], DebugLevel); err != nil {
		fmt.Println("Set logger Error:", err)
		return
	}
	defer Common.CloseLogger()

	//获取网络配置
	HttpConf, err := conf.GetConf("HTTP_CONF")
	if err != nil {
		Common.FATAL("get configer faile. Error:", err)
	}

	//启动Web Server
	http.HandleFunc(HttpConf["ListenHos"], DisplyPck.DisplyIndex)

}