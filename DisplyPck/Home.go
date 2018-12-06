package DisplyPck

import (
	"html/template"
	"net/http"
)

type IndexInfo struct {
	Title   string
	Date    string
	Weather string
}

func DisplyIndex(w http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("./HTMLtemplate/IndexTemplate.html"))
	indexValue := IndexInfo{Title: "Test Page", Date: "2018-12-5", Weather: "Fine"}
	tmpl.Execute(w, indexValue)
}
