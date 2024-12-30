package controller

import (
	"fmt"
	"net/http"

	"github.com/mjishu/pokeDate/database"
)

func UserController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" {
		handleUsers(w, r)
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
