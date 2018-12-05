package main

import (
	"Common"
	"flag"
	"net/https"
)

var ConfFile = flag.String("f", "./WevServer.conf", "configer file")

func main() {
	//解析配置
	conf := &Common.Configer{}
	if err := conf.Init(ConfFile); err != nil {

	}
}
