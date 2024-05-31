package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Create a new ServeMux using Gorilla
var rMux = mux.NewRouter()

// PORT is where the web server listens to
var PORT = ":1234"


func main () {
	arguments := os.Args

	if len(arguments) >= 2 {
		PORT = ":" + arguments[1]
	}
	
	s := http.Server{
		Addr: PORT,
		Handler: rMux,
		ErrorLog: nil,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout: 10 * time.Second,
	}

	rMux.NotFoundHandler = http.HandlerFunc(DefaultHandler)

	notAllowed := notAllowedHandler{}
	rMux.MethodNotAllowedHandler = notAllowed

	rMux.HandleFunc("/time", TimeHandler)
	








}