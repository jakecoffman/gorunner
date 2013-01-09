package server

import (
	"net/http"
	"fmt"
)

func Handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello2");
}
