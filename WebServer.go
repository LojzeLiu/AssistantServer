package main

import (
	"Common"
	"flag"
	"fmt"
	"net/https"
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

	//启动Web Server
}
