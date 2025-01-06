package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mjishu/pokeDate/database"
)

func ShotController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/shots" {
		fmt.Fprintf(w, "Incorrect Path")
		return
	}
	SetHeader(w)

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(database.GetAllShots()); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
	}
}

func OrganizationController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/organizations/animals" {
		fmt.Fprintf(w, "Incorrect Path")
		return
	}

	SetHeader(w)

	switch r.Method {
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")

		hasId, id := checkForBodyItem("id", w, r)
		animal, err := database.GetAnimal(id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "could not find animal", err)
			return
		}
		if hasId {
			if err := json.NewEncoder(w).Encode(animal); err != nil {
				http.Error(w, "unable to encode response", http.StatusInternalServerError)
			}
			return
		}
		fmt.Fprintf(w, "Body does not have an id!")

	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(database.GetAllAnimals()); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
		return
	}
}
