package DisplyPck

import (
	"AssistantServer/NetLayer"
	"html/template"
	"net/http"
)

func DisplyIndex(w http.ResponseWriter, req *http.Request) {
	TodayWeather, TodayEarly, err := gCrawler.GetTodayBrief(101010100)
	if err != nil {
		errTmp := template.Must(template.ParseFiles("./tmpl"))
		errData := ErrorInfo{Msg: err}
		errTmp.Execute(w, errData)
		return
	}
	tmpl := template.Must(template.ParseFiles("./tmpl/IndexTemplate.html"))
	tmpl.Execute(w, TodayWeather)
}
