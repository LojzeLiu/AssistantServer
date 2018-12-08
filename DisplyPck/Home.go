package DisplyPck

import (
	"Common"
	"fmt"
	"html/template"
	"net/http"
)

func DisplyIndex(w http.ResponseWriter, req *http.Request) {
	TodayWeather, TodayEarly, err := gCrawler.GetTodayBrief(2)
	if err != nil {
		errTmp := template.Must(template.ParseFiles("./tmpl/ErrorTemp.html"))
		errData := ErrorInfo{Title: "Error", Msg: fmt.Sprint("Get Today Brief Failed. Reason:", err)}
		errTmp.Execute(w, errData)
		Common.ERROR("Get Today Brief Failed. Reason:", err)
		return
	}
	if TodayEarly != nil {
		Common.DEBUG("Early Weather")
	}
	TodayWeather.Title = "Today Weather"
	tmpl := template.Must(template.ParseFiles("./tmpl/IndexTemplate.html"))
	tmpl.Execute(w, TodayWeather)
}
