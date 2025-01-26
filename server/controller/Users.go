package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

type AuthUser struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func UserController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket, s3Region string) {
	SetHeader(w)
	switch r.URL.Path {
	case "/users/login":
		LoginUser(w, r, pool, jwtSecret)
		return
	case "/users/create":
		switch r.Method {
		case http.MethodPost:
			CreateUser(w, r, pool, jwtSecret, s3Bucket, s3Region)
		case http.MethodGet:
			fmt.Fprint(w, "no get route setup")
		}
		return
	case "/users/current":
		switch r.Method {
		case http.MethodPost:
			GetCurrentUser(w, r.Header, pool, jwtSecret, s3Bucket, s3Region)
		case http.MethodPut:
			UpdateUser(w, r, pool, jwtSecret)
		}
		return
	default:
		fmt.Println("This is the default path")
		return
	}
}

func LoginUser(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) { //? does this properly check if the usernames are the same before logging in?
	var incomingUser AuthUser
	var expiresIn time.Duration
	checkAuthUser(w, r, &incomingUser) //* gets req body

	storedUser, err := database.GetUser(pool, incomingUser.Username) //this password should be hashed(i.e user.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find user", err)
	}

	//* checks password. useful code goes after this
	err = auth.CheckPasswordHash(incomingUser.Password, storedUser.HashPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid password", err)
		return
	}

	//* Token information
	expiresIn = time.Duration(15 * time.Minute)

	token, err := auth.MakeJWT(storedUser.Id, jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create JWT", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create refresh token", err)
		return
	}
	_, err = database.CreateRefreshToken(pool, refresh_token, storedUser.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not store refresh token", err)
		return
	}

	response := map[string]interface{}{
		"username":      storedUser.Username,
		"id":            storedUser.Id,
		"status":        http.StatusOK,
		"token":         token,
		"refresh_token": refresh_token,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func CreateUser(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket, s3Region string) {
	var user database.User
	checkBody(w, r, &user)
	fmt.Printf("user in body is %v\n", user) // this is responding with a nil uuid

	hashedPassword, err := auth.HashPassword(user.HashPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not hash password", err)
		return
	}

	defaultPfp := "https://" + s3Bucket + ".s3." + s3Region + ".amazonaws.com/profile_pictures/default.webp"
	user.Profile_picture = &defaultPfp

	// creates user
	user.HashPassword = hashedPassword
	storedId, err := database.CreateUser(pool, user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create user", err)
		return
	}

	//* Token information
	expiresIn := time.Duration(15 * time.Minute)

	token, err := auth.MakeJWT(storedId, jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create JWT", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create refresh token", err)
		return
	}
	_, err = database.CreateRefreshToken(pool, refresh_token, storedId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not store refresh token", err)
		return
	}

	response := map[string]interface{}{
		"username":      user.Username,
		"id":            user.Id,
		"status":        http.StatusOK,
		"token":         token,
		"refresh_token": refresh_token,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func GetCurrentUser(w http.ResponseWriter, header http.Header, pool *pgxpool.Pool, jwtSecret, s3Bucket, s3Region string) {
	tokenUserId, err := auth.UserValid(header, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find jwt", err)
		return
	}
	storedUser, err := database.GetUserById(pool, tokenUserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find user by id", err)
		return
	}
	respondWithJSON(w, http.StatusOK, storedUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	fmt.Println("put was called")
	userIdToken, err := auth.UserValid(r.Header, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find JWT", err)
		return
	}

	userData, err := database.GetUserById(pool, userIdToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find user with id", err)
		return
	}

	var incomingUser database.User
	checkUpdateUser(w, r, &incomingUser)

	if userData.Username != incomingUser.Username {
		userData.Username = incomingUser.Username
	}
	if userData.Email != incomingUser.Email {
		userData.Email = incomingUser.Email
	}
	if userData.Date_of_birth != incomingUser.Date_of_birth {
		userData.Date_of_birth = incomingUser.Date_of_birth
	}

	fmt.Printf("updated user is %v\n", userData)

	err = database.UpdateUser(pool, userData)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not update user", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

//? -------------------- GETS item from body

func checkBody(w http.ResponseWriter, r *http.Request, user interface{}) error {
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
	fmt.Printf("checkBody info is %v\n", user)
	return err
}

// Modify the existing functions to use the new checkUser function
func checkAuthUser(w http.ResponseWriter, r *http.Request, user *AuthUser) error {
	return checkBody(w, r, user)
}

func checkUpdateUser(w http.ResponseWriter, r *http.Request, user *database.User) error {
	return checkBody(w, r, user)
}
