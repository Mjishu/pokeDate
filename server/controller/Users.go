package controller

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/mjishu/pokeDate/database"
)

func UserController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" {
		handleUsers(w, r)
	} else if r.URL.Path == "/users/login" {
		LoginUser(w, r)
	} else if r.URL.Path == "/users/signup" {
		CreateUser(w, r)
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) database.User { //? how to send the user back to frontend?
	switch r.Method {
	case http.MethodGet:
		hasid, userId := checkForBodyItem("id", w, r)
		if !hasid {
			fmt.Fprint(w, "Could not find the user id in body", http.StatusBadRequest)
		}
		user, err := database.GetUser(userId)
		if err != nil {
			fmt.Fprint(w, "There was an error trying to create user", http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "user recieved", http.StatusAccepted)
		return user
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) { //? how to get this to work so that it passes the user of body to createUser

	_, user := checkForBodyItem("User", w, r)

	// _, ok := user.(database.NewUser) // checks if user is of type NewUser
	// if !ok {
	// 	return
	// }
	if reflect.TypeOf(user) == reflect.TypeOf(database.NewUser{}) {
		database.CreateUser(user)
		return
	}
}
