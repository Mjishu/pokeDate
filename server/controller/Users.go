package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

type AuthUser struct {
	Username string
	Password string
}

func UserController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" {
		handleUsers(w, r)
	} else if r.URL.Path == "/users/login" {
		LoginUser(w, r)
	} else if r.URL.Path == "/users/create" {
		switch r.Method {
		case http.MethodPost:
			CreateUser(w, r)
		case http.MethodGet:
			fmt.Fprint(w, "no get route setup")
		}
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
		fmt.Fprint(w, "user recieved", http.StatusAccepted)
		return user
	}
	return database.User{}
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var incomingUser AuthUser
	checkAuthUser(w, r, &incomingUser)

	storedUser, err := database.GetUser(incomingUser.Username) //this password should be hashed(i.e user.Password)
	if err != nil {
		fmt.Fprint(w, "Error getting user from database", http.StatusInternalServerError)
	}

	err = auth.CheckPasswordHash(incomingUser.Password, storedUser.HashPassword)
	if err != nil {
		fmt.Fprint(w, "issue checking passwords", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"username": storedUser.Username,
		"id":       storedUser.Id,
		"status":   http.StatusOK,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) { //? how to get this to work so that it passes the user of body to createUser

	var user database.NewUser
	checkUser(w, r, &user)

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		fmt.Fprint(w, "error trying to hash password", http.StatusInternalServerError)
		return
	}
	database.CreateUser(user, hashedPassword)
	fmt.Printf("hash is %s\n", hashedPassword)
}

//? -------------------- GETS item from body

func checkAuthUser(w http.ResponseWriter, r *http.Request, user *AuthUser) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return nil
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusInternalServerError)
	}
	return err
}

func checkUser(w http.ResponseWriter, r *http.Request, user *database.NewUser) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return nil
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusInternalServerError)
	}
	return err
}
