package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Job(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "<html>Hello job: %s!</html>", vars["job"])
}
