package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mjishu/pokeDate/database"
)

func AnimalController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.URL.Path == "/organizations/animals" {
		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")

		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(database.GetAllAnimals()); err != nil {
				http.Error(w, "unable to encode response", http.StatusInternalServerError)
			}
			return
		}
	}

	if r.URL.Path == "/animals" {

		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")

			animal := GetAnimalFromBody("animal", w, r)

			database.InsertAnimal(animal)
		}
	}
}

func GetFrontendURL() string {
	err := godotenv.Load("../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal("ERror loading .env file")
		}
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatal("DATABASE_URL not set in .env")
	}
	return frontendURL
}

func GetAnimalFromBody(key string, w http.ResponseWriter, r *http.Request) database.NewAnimal {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return database.NewAnimal{}
	}
	defer r.Body.Close()

	var animal database.NewAnimal
	err = json.Unmarshal(body, &animal)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return database.NewAnimal{}
	}
	return animal
}
