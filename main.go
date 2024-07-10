package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// Create a new ServeMux using Gorilla
var rMux = mux.NewRouter()

// PORT is where the web server listens to
var PORT = ":1234"

// SliceToJSON encodes a slice with JSON records

func main() {
	arguments := os.Args

	if len(arguments) >= 2 {
		PORT = ":" + arguments[1]
	}

	s := http.Server{
		Addr:         PORT,
		Handler:      rMux,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	rMux.NotFoundHandler = http.HandlerFunc(DefaultHandler)

	notAllowed := notAllowedHandler{}
	rMux.MethodNotAllowedHandler = notAllowed

	rMux.HandleFunc("/time", TimeHandler)

	getMux := rMux.Methods(http.MethodGet).Subrouter()
	getMux.HandleFunc("/getall", GetAllHandler)
	getMux.HandleFunc("/getid/{username}", GetIdHandler)
	getMux.HandleFunc("/logged", LoggedUsersHandler)
	getMux.HandleFunc("/username/{id:[0-9]+}", GetUserDataHandler)

	putMux := rMux.Methods(http.MethodPut).Subrouter()
	putMux.HandleFunc("/update", UpdateHandler)

	postMux := rMux.Methods(http.MethodPost).Subrouter()
	postMux.HandleFunc("/add", UpdateHandler)
	postMux.HandleFunc("/login", LoginHandler)
	postMux.HandleFunc("/logout", LogoutHandler)

	deleteMux := rMux.Methods(http.MethodDelete).Subrouter()
	deleteMux.HandleFunc("/username/{id:[0-9]+}", DeleteHandler)


	go func ()  {
		log.Println("Listening to", PORT)
		err := s.ListenAndServe()

		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			return
		}
	}()


	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <- sigs
	log.Println("Quitting after signal:", sig)
	time.Sleep(5 * time.Second)
	s.Shutdown(nil)
	
}
