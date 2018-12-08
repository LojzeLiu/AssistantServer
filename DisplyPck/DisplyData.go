package DisplyPck

import (
	"AssistantServer/NetLayer"
	"Common"
)

var gCrawler NetLayer.WeatherCrawler

type DisplyData struct {
	mConf *Common.Configer
}

func (this *DisplyData) Init(conf *Common.Configer) error {
	this.mConf = conf
	return gCrawler.Init(conf) //初始天气爬虫
}

type ErrorInfo struct {
	Title string
	Msg   string
}
