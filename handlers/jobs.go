package handlers

import (
	"fmt"
	"net/http"
)

func Jobs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html>Hello jobs!</html>")
}
