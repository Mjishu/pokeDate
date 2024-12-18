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
	SetHeader(w)

	if r.URL.Path == "/animals" {

		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			fmt.Println("The else in method post was called")

			animal := GetAnimalFromBody(w, r)

			database.InsertAnimal(animal)
			fmt.Fprintf(w, "Animal created Successfully!: %v\n", animal)
		case http.MethodPut:
			w.Header().Set("Content-Type", "application/json")

			updatedAnimal := GetUpdatedAnimalFromBody(w, r)

			database.UpdateAnimal(updatedAnimal)
			fmt.Fprintf(w, "Animal updated successfully")
		case http.MethodDelete:
			w.Header().Set("Content-Type", "application/json")

			hasId, id := checkForBodyItem("id", w, r)

			if hasId {
				database.DeleteAnimal(id)
				fmt.Fprintf(w, "Animal was removed successfully")
				return
			}
			fmt.Fprintf(w, "Id was not found in body!")
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

func GetAnimalFromBody(w http.ResponseWriter, r *http.Request) database.NewAnimal {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error trying to read body: %v\n", err)
		http.Error(w, "unable to read body: %v\n", http.StatusInternalServerError)
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
