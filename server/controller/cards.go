package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
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

func CardsController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	if r.URL.Path != "/cards" {
		http.NotFound(w, r)
		return
	}

	SetHeader(w)
	switch r.Method {

	//* GET
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")

		animal := database.GetRandomAnimal(pool)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(animal); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
		return
	//* POST
	case http.MethodPost:

	}
}

func UserCardResponse(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	/*todo
	create message request with the organization who has this animal andour userId

	*/
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	_, err = auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	SetHeader(w)
	var card Card
	err = checkBody(w, r, &card)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find card id in body", err)
		return
	}

	animal, err := database.GetAnimal(pool, card.Id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animal", err)
		return
	}

	respondWithJSON(w, http.StatusOK, animal)
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
