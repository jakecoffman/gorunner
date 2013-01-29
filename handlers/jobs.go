package handlers

import (
	"html/template"
	"net/http"
)

var index = template.Must(template.ParseFiles(
	"web/templates/_base.html",
	"web/templates/index.html",
))

func Jobs(w http.ResponseWriter, r *http.Request) {
	if err := index.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
