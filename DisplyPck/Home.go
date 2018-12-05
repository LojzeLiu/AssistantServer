package DisplyPck

import (
	"Common"
	"html/template"
	"net/http"
)

type IndexInfo struct {
	Title   string
	Date    string
	Weather string
}

func DisplyIndex(w http.ResponseWriter, req *http.Request) {
	IndexTemp := template.New("Index")
	IndexTemp, err := IndexTemp.ParseFiles("./HTMLtemplate/IndexTemplate.html")
	if err != nil {
		Common.ERROR("Parse File Error:", err)
		return
	}

	var index IndexInfo
	index.Title = "Test Page"
	index.Date = "2018-12-5"
	index.Weather = "晴"
	IndexTemp.Execute(w, index)
}
