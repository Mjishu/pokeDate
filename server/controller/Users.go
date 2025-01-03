package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

type AuthUser struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func UserController(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	SetHeader(w)
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/users/login":
		LoginUser(w, r, jwtSecret)
		return
	case "/users/create":
		switch r.Method {
		case http.MethodPost:
			CreateUser(w, r)
		case http.MethodGet:
			fmt.Fprint(w, "no get route setup")
		}
		return
	case "/users/current":
		GetCurrentUser(w, r, jwtSecret)
		return
	default:
		handleUsers(w, r)
		return
	}
	// if r.URL.Path == "/users" {
	// 	handleUsers(w, r)
	// 	return
	// } else if r.URL.Path == "/users/login" {
	// 	LoginUser(w, r, jwtSecret)
	// 	return
	// } else if r.URL.Path == "/users/create" {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		CreateUser(w, r)
	// 	case http.MethodGet:
	// 		fmt.Fprint(w, "no get route setup")
	// 	}
	// 	return
	// }
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

func LoginUser(w http.ResponseWriter, r *http.Request, jwtSecret string) { //? does this properly check if the usernames are the same before logging in?
	var incomingUser AuthUser
	var expiresIn time.Duration
	checkAuthUser(w, r, &incomingUser) //* gets req body

	fmt.Printf("incoming user is %v\n", incomingUser)
	storedUser, err := database.GetUser(incomingUser.Username) //this password should be hashed(i.e user.Password)
	if err != nil {
		fmt.Fprint(w, "Error getting user from database", http.StatusInternalServerError)
	}

	//* checks password. useful code goes after this
	err = auth.CheckPasswordHash(incomingUser.Password, storedUser.HashPassword)
	if err != nil {
		fmt.Fprint(w, "issue checking passwords", http.StatusBadRequest)
		return
	}

	// Token information
	expiresIn = time.Duration(15 * time.Minute)

	token, err := auth.MakeJWT(storedUser.Id, jwtSecret, expiresIn)
	if err != nil {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"success": false}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		fmt.Printf("error finding json token %v\n ", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"success": false}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		fmt.Printf("error finding getting refresh token %v\n ", err)
		return
	}
	_, err = database.CreateRefreshToken(refresh_token, storedUser.Id)
	if err != nil {
		http.Error(w, "error creating refresh token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"username":      storedUser.Username,
		"id":            storedUser.Id,
		"status":        http.StatusOK,
		"token":         token,
		"refresh_token": refresh_token,
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

func GetCurrentUser(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	switch r.Method {
	case http.MethodPost:
		err := auth.UserValid(r.Header, jwtSecret)
		if err != nil {
			http.Error(w, "unable to validate jwt", http.StatusBadRequest)
			return
		}
		/*
			get current user id from jwt, then search up that user in db and return that user
		*/

		response := map[string]interface{}{
			"username": "user",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
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
