package handlers

import (
	"html/template"
	"net/http"
)

const base string = "github.com/jakecoffman/gorunner/web/"

var index = template.Must(template.ParseFiles(
	base+"templates/_base.html",
	base+"templates/index.html",
))

func Jobs(w http.ResponseWriter, r *http.Request) {
	if err := index.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
