package controller

import (
	"encoding/json"
	"fmt"
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
			hasId, id := checkForBodyItem("id", w, r)
			if hasId {
				fmt.Println("Id was found in body")
				fmt.Printf("id is %s", id)
				if err := json.NewEncoder(w).Encode(database.GetAnimal(id)); err != nil {
					http.Error(w, "unable to encode response", http.StatusInternalServerError)
				}
				return

			} else {
				w.Header().Set("Content-Type", "application/json")

				animal := GetAnimalFromBody(w, r)

				database.InsertAnimal(animal)
				fmt.Fprintf(w, "Animal created Successfully!")
			}
		case http.MethodPut:
			fmt.Println("put has been called")
			w.Header().Set("Content-Type", "application/json")

			updatedAnimal := GetUpdatedAnimalFromBody(w, r)

			database.UpdateAnimal(updatedAnimal)
			fmt.Fprintf(w, "Animal updated successfully")
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

func GetAnimalFromBody(w http.ResponseWriter, r *http.Request) database.NewAnimal { // *can i make this take a type as a parmater and change the animal database.NewAnimal based on that?
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
	fmt.Printf("Animal inside GAFB is %v\n", animal)
	return animal
}

func GetUpdatedAnimalFromBody(w http.ResponseWriter, r *http.Request) database.UpdateAnimalStruct {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return database.UpdateAnimalStruct{}
	}
	defer r.Body.Close()

	var animal database.UpdateAnimalStruct
	err = json.Unmarshal(body, &animal)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return database.UpdateAnimalStruct{}
	}
	fmt.Printf("Animal inside GAFB is %v\n", animal)
	return animal
}

func checkForBodyItem(key string, w http.ResponseWriter, r *http.Request) (bool, any) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return false, ""
	}
	defer r.Body.Close()

	var elements map[string]interface{}
	err = json.Unmarshal(body, &elements)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusInternalServerError)
		return false, ""
	}

	if _, exists := elements[key]; exists {
		return true, elements[key]
	}
	return false, ""
}
