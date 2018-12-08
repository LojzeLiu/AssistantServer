package DisplyPck

import (
	"AssistantServer/NetLayer"
	"Common"
	"fmt"
	"html/template"
	"net/http"
)

type AlertShow struct {
	AlertLeve string
	IsAlert   bool
	Weather   *NetLayer.TodayWeatherBrief
	Alert     []NetLayer.TodayAlertWeather
}

func (this AlertShow) String() string {
	return fmt.Sprintf("Weather:%s;Alert:%s", this.Weather, this.Alert)
}

func DisplyIndex(w http.ResponseWriter, req *http.Request) {
	TodayWeather, TodayEarly, err := gCrawler.GetTodayBrief(2)
	if err != nil {
		errTmp := template.Must(template.ParseFiles("./tmpl/ErrorTemp.html"))
		errData := ErrorInfo{Title: "Error", Msg: fmt.Sprint("Get Today Brief Failed. Reason:", err)}
		errTmp.Execute(w, errData)
		Common.ERROR("Get Today Brief Failed. Reason:", err)
		return
	}
	Show := AlertShow{}
	TodayWeather.Title = "Today Weather"
	Show.Weather = TodayWeather
	if TodayEarly != nil {
		Common.DEBUG("Early Weather")
		Show.IsAlert = true
		Show.Alert = TodayEarly
		Show.AlertLeve = "#FF0000"
	}
	tmpl := template.Must(template.ParseFiles("./tmpl/IndexTemplate.html"))
	tmpl.Execute(w, Show)
}
