package main

import (
	"AssistantServer/DisplyPck"
	"AssistantServer/NetLayer"
	"Common"
	"flag"
	"golang.org/x/net/websocket"
	"net/http"
	"strconv"
)

var ConfFile = flag.String("f", "./WebServer.conf", "configer file")

func main() {
	//解析配置
	conf := &Common.Configer{}
	if err := conf.Init(*ConfFile); err != nil {
		return
	}
	defer conf.Destroy()

	//获取日志配置
	logConfs, err := conf.GetConf("LOG")
	if err != nil {
		return
	}
	DebugLevel, err := strconv.Atoi(logConfs["Level"])
	if err != nil {
		return
	}

	//初始化日志
	if err := Common.SetLogger(logConfs["Path"], logConfs["AppName"], Common.LOG_LEVE(DebugLevel)); err != nil {
		return
	}
	defer Common.CloseLogger()

	//获取网络配置
	HttpConf, err := conf.GetConf("HTTP_CONF")
	if err != nil {
		Common.FATAL("get configer faile. Error:", err)
	}
	if len(HttpConf) <= 0 {
		Common.WARN("Not Http configer.")
	}

	//配置静态文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//配置首页
	http.HandleFunc("/", DisplyPck.DisplyIndex)

	//配置Websoekct
	http.Handle("/ws_accept/", websocket.Handler(NetLayer.WSAccept))

	//启动Web Server
	if err := http.ListenAndServe(HttpConf["ListenHos"], nil); err != nil {
		Common.ERROR("Error:", err)
		return
	}

}
