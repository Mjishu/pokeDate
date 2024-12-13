package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mjishu/pokeDate/database"
)

type Animal struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Species       string `json:"species"`
	Date_of_birth string `json:"date_of_birth"`
	Sex           string `json:"sex"`
	Available     bool   `json:"available"`
	Image_src     string `json:"image_src"`
}

type Card struct {
	Id              string `json:"id"`
	Animal_id       string `josn:"animal_id"`
	Organization_id string `json:"organization_id"`
	Liked           bool   `json:"liked"`
	Animal_info     Animal `json:"animal_info"`
}

func CardsController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cards" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", GetFrontendURL())
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	switch r.Method {

	//* GET
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		animal := database.GetRandomAnimal()

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(animal); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
		return
	//* POST
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		id := GetIdFromBody("id", w, r)

		animal := database.GetAnimals(id)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(animal); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
		return
	}
	GetFromHeader("Authorization", r)
}

func GetFromHeader(key string, r *http.Request) string {
	for name, values := range r.Header {
		if name == key {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
				return value
			}
		}
	}
	return ""
}

func GetIdFromBody(key string, w http.ResponseWriter, r *http.Request) string {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return ""
	}
	defer r.Body.Close()

	var card Card
	err = json.Unmarshal(body, &card)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return ""
	}
	return card.Id
}
