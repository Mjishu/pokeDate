package cards

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type Card struct {
	Id              string `json:"id"`
	Animal_id       string `josn:"animal_id"`
	Organization_id string `json:"organization_id"`
	Liked           bool   `json:"liked"`
	Animal_info     Animal `json:"animal_info"`
}

type Animal struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Species       string `json:"species"`
	Date_of_birth string `json:"date_of_birth"`
	Sex           string `json:"sex"`
	Available     bool   `json:"available"`
	Image_src     string `json:"image_src"`
}

var animalData = []Animal{
	{Id: "001", Name: "Nilah", Species: "dog", Date_of_birth: "2020/04/23", Sex: "female", Available: true, Image_src: "./images/dog.webp"},
	{Id: "002", Name: "fortuna", Species: "cat", Date_of_birth: "2023/09/02", Sex: "male", Available: false, Image_src: "./images/cat.jpg"},
}

var data = []Card{
	{Id: "001", Animal_id: "001", Organization_id: "001", Animal_info: animalData[0]},
	{Id: "002", Animal_id: "002", Organization_id: "002", Animal_info: animalData[1]},
	{Id: "003", Animal_id: "003", Organization_id: "003", Animal_info: animalData[0]},
}

func CardsController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cards" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	switch r.Method {

	//* GET
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		random := rand.Intn(len(data))
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data[random]); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
		return
	//* POST
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		id := GetIdFromBody("id", w, r)

		for _, card := range data {
			if card.Id == id {
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(card); err != nil {
					http.Error(w, "unable to encode response", http.StatusInternalServerError)
				}
				return
			}
		}
	}
	fmt.Println("looking for cards to control")
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
