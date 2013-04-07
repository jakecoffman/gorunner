package handlers

import (
	"html/template"
	"net/http"
	"github.com/jakecoffman/gorunner/models"
	"github.com/jakecoffman/gorunner/utils"
	"github.com/gorilla/mux"
)

func Triggers(w http.ResponseWriter, r *http.Request) {
	triggerList := models.GetTriggerList()

	if r.Method == "GET" {
		t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
			"web/templates/_base.html",
			"web/templates/_nav.html",
			"web/templates/triggers.html",
		))

		if err := t.Execute(w, triggerList); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "POST" {
		trigger := r.FormValue("name")
		triggerList.Append(trigger)
	} else {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func Trigger(w http.ResponseWriter, r *http.Request){
	triggerList := models.GetTriggerList()

	vars := mux.Vars(r)
	trigger, err := triggerList.Get(vars["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	t := template.Must(template.New("_base.html").Funcs(utils.FuncMap).ParseFiles(
		"web/templates/_base.html",
		"web/templates/_nav.html",
		"web/templates/trigger.html",
	))

	if err := t.Execute(w, trigger); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
