package main

import (
	"fmt"
	"net/http"
	"github.com/jakecoffman/gorunner/server"
)

func main() {
	http.HandleFunc("/", server.Handler)
	http.HandleFunc("/(.*)", server.Handler2)
	fmt.Println("Running on port 8090")
	http.ListenAndServe(":8090", nil)
}