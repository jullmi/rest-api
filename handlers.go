package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	restdb "github.com/jullmi/restdb"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"user"`
	Password  string `json:"password"`
	LastLogin int64  `json:"lastlogin"`
	Admin     int    `json:"admin"`
	Active    int    `json:"active"`
}

// SliceToJSON encodes a slice with JSON records
func SliceToJSON(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(slice)
}

type notAllowedHandler struct{}

func (h notAllowedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	MethodNotAllowedHandler(rw, r)
}

// MethodNotAllowedHandler is executed when url path not found
func DefaultHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("DefaultHandler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)

	rw.WriteHeader(http.StatusNotFound)
	Body := r.URL.Path + " is not supported. Thanks for visiting!\n"
	fmt.Fprintf(rw, "%s", Body)

}

// MethodNotAllowedHandler is executed when the HTTP method is incorrect
func MethodNotAllowedHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)

	rw.WriteHeader(http.StatusNotFound)
	Body := "Method not allowed!\n"
	fmt.Fprintf(rw, "%s", Body)
}

// TimeHandler is for handling /time – it works with plain text
func TimeHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("TimeHandler Serving:", r.URL.Path, "from", r.Host)

	rw.WriteHeader(http.StatusOK)
	t := time.Now().Format(time.RFC1123)

	Body := "The current time is: " + t + "\n"
	fmt.Fprint(rw, "s", Body)

}

// AddHandler is for adding a new user
func AddHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("AddHandler Serving:", r.URL.Path, "from", r.Host)

	d, err := io.ReadAll(r.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	if len(d) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	// We read two structures as an array:
	// 1. The user issuing the command
	// 2. The user to be added

	var users = []restdb.User{}
	err = json.Unmarshal(d, &users)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	log.Println(users)

	if !restdb.IsUserAdmin(users[0]) {
		log.Println("Command issued by non-admin user:", users[0].Username)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	result := restdb.InsertUser(users[1])
	if !result {
		rw.WriteHeader(http.StatusBadRequest)
	}

}

func DeleteHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("DeleteHandler Serving:", r.URL.Path, "from", r.Host)

	// Get the ID of the user to be deleted
	id, ok := mux.Vars(r)["id"]

	if !ok {
		log.Println("ID value not set!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := restdb.User{}
	err := user.FromJSON(r.Body)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !restdb.IsUserAdmin(user) {
		log.Println("User", user.Username, "is not admin!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(id)

	if err != nil {
		log.Println("id", err)
		return
	}

	t := restdb.FindUserID(intID)

	if t.Username != "" {
		log.Println("About to delete:", t)
		deleted := restdb.DeleteUser(intID)

		if deleted {
			log.Println("User deleted:", id)
			rw.WriteHeader(http.StatusOK)

		} else {
			log.Println("User ID not found:", id)
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}

	rw.WriteHeader(http.StatusNotFound)
}

// GetAllHandler is for getting all data from the user database
func GetAllHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("GetAllHandler Serving:", r.URL.Path, "from", r.Host)

	d, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(d) == 0 {
		log.Println("No input!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := restdb.User{}
	err = json.Unmarshal(d, &user)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !restdb.IsUserAdmin(user) {
		log.Println("User", user, "is not an admin!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = SliceToJSON(restdb.ListAllUsers(), rw)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

// GetIDHandler returns the ID of an existing user
func GetIdHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("GetIDHandler Serving:", r.URL.Path, "from", r.Host)

	username, ok := mux.Vars(r)["username"]

	if !ok {
		log.Println("ID value not set!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	d, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(d) == 0 {
		log.Println("No input!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := restdb.User{}
	err = json.Unmarshal(d, &user)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Input user:", user)

	if !restdb.IsUserAdmin(user) {
		log.Println("User", user.Username, "not an admin!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	t := restdb.FindUserUsername(username)
	if t.ID != 0 {
		err = t.ToJSON(rw)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("User " + user.Username + "not found")
	}
}

// LoggedUsersHandler returns the list of all logged in users
func LoggedUsersHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("LoggedUsersHandler Serving:", r.URL.Path, "from", r.Host)

	user := restdb.User{}
	err := user.FromJSON(r.Body)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !restdb.IsUserValid(user) {
		log.Println("User", user.Username, "exists!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = SliceToJSON(restdb.ReturnLoggedUsers(), rw)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

}

// GetUserDataHandler + GET returns the full record of a user
func GetUserDataHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("GetUserDataHandler Serving:", r.URL.Path, "from", r.Host)

	id, ok := mux.Vars(r)["id"]

	if !ok {
		log.Println("ID value not set!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	intId, err := strconv.Atoi(id)

	if err != nil {
		log.Println("id", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	t := restdb.FindUserID(intId)

	if t.ID != 0 {
		err = t.ToJSON(rw)

		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

// UpdateHandler is for updating the data of an existing user + PUT
func UpdateHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("UpdateHandler Serving:", r.URL.Path, "from", r.Host)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(d) == 0 {
		log.Println("No input!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	users := []restdb.User{}
	err = json.Unmarshal(d, &users)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !restdb.IsUserAdmin(users[0]) {
		log.Println("Command issued by non-admin user:", users[0].Username)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(users)

	t := restdb.FindUserUsername(users[1].Username)
	t.Username = users[1].Username
	t.Admin = users[1].Admin
	t.Password = users[1].Password

	if !restdb.UpdateUser(t) {
		log.Println("Update failed:", t)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Update successful:", t)
	rw.WriteHeader(http.StatusOK)
}

// LoginHandler is for updating the LastLogin time of a user
// And changing the Active field to true
func LoginHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler Serving:", r.URL.Path, "from", r.Host)

	d, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(d) == 0 {
		log.Println("No input!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := restdb.User{}

	err = json.Unmarshal(d, &user)
	if err !=nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Input user:", user)

	if !restdb.IsUserValid(user) {
		log.Println("User", user.Username, "not valid!")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}



}
