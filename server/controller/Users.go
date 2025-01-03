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
	Username           string `json:"Username"`
	Password           string `json:"Password"`
	Expires_in_seconds int    `json:"exp_seconds"`
}

func UserController(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	SetHeader(w)
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/users" {
		handleUsers(w, r)
		return
	} else if r.URL.Path == "/users/login" {
		LoginUser(w, r, jwtSecret)
		return
	} else if r.URL.Path == "/users/create" {
		switch r.Method {
		case http.MethodPost:
			CreateUser(w, r)
		case http.MethodGet:
			fmt.Fprint(w, "no get route setup")
		}
		return
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

func LoginUser(w http.ResponseWriter, r *http.Request, jwtSecret string) { //? does this properly check if the usernames are the same before logging in?
	var incomingUser AuthUser
	var expiresIn time.Duration
	maxExpireTime := time.Duration(1 * time.Hour)
	checkAuthUser(w, r, &incomingUser)

	fmt.Printf("incoming user is %v\n", incomingUser)
	storedUser, err := database.GetUser(incomingUser.Username) //this password should be hashed(i.e user.Password)
	if err != nil {
		fmt.Fprint(w, "Error getting user from database", http.StatusInternalServerError)
	}

	err = auth.CheckPasswordHash(incomingUser.Password, storedUser.HashPassword)
	if err != nil {
		fmt.Fprint(w, "issue checking passwords", http.StatusBadRequest)
		return
	}

	// checks if time is > than max expire time if so set expiresIn to maxTime
	expiresTime := time.Duration(incomingUser.Expires_in_seconds) * time.Second
	if expiresTime == 0 {
		expiresIn = maxExpireTime // seconds * minutes = 1 hour
	} else if expiresTime > maxExpireTime { // shuld reject anything past 3600
		expiresIn = maxExpireTime
	} else {
		expiresIn = expiresTime
	}
	fmt.Printf("Expires in: %v\n", expiresIn)
	// need to call the token somewhere here

	token, err := auth.MakeJWT(storedUser.Id, jwtSecret, expiresIn)
	if err != nil {
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"success": false}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		fmt.Printf("error finding json token %v", err)
		return
	}

	response := map[string]interface{}{
		"username": storedUser.Username,
		"id":       storedUser.Id,
		"status":   http.StatusOK,
		"token":    token,
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
