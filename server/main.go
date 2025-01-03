package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/controller"
	"github.com/mjishu/pokeDate/database"
)

type apiConfig struct {
	jwt_secret string
}

func main() {
	mux := http.NewServeMux()
	cld, ctx := database.Credentials()
	config := &apiConfig{}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.loadParams()

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		controller.UserController(w, r, config.jwt_secret)
	})
	mux.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		err := auth.UserValid(r.Header, config.jwt_secret)
		if err != nil {
			http.Error(w, "unable to validate jwt", http.StatusBadRequest)
			return
		}

		//* CONTROLLER
		controller.CardsController(w, r)
	})
	mux.HandleFunc("/animals/", func(w http.ResponseWriter, r *http.Request) {
		controller.AnimalController(w, r, cld, ctx)
	})
	mux.HandleFunc("/organizations/animals", controller.OrganizationController) //? change to /orgnaizations and make a new controller called organizations Cotnroller
	mux.HandleFunc("/shots", controller.ShotController)
	mux.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		controller.RefreshToken(w, r, config.jwt_secret)
	})
	mux.HandleFunc("/revoke", controller.RevokeToken)

	database.Database()

	port := ":8080"

	fmt.Println("listening on port " + port)
	err = http.ListenAndServe(port, mux)
	log.Fatal(err)
}

func (config *apiConfig) loadParams() {
	config.jwt_secret = os.Getenv("JWT_SECRET")
}
