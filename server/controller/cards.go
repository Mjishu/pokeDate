package controller

import (
	"encoding/json"
	"fmt"
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
	Animal_id string `josn:"Animal_id"`
	Liked     bool   `json:"Liked"`
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
		GetRandomCard(w, r, pool, jwtSecret)
	//* POST
	case http.MethodPost:
		UserCardResponse(w, r, pool, jwtSecret)
	}
}

func GetRandomCard(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	userId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	animal, err := database.GetRandomAnimal(pool, userId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get random animal", err)
		return
	}

	respondWithJSON(w, http.StatusOK, animal)
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

	userId, err := auth.ValidateJWT(token, jwtSecret)
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

	fmt.Printf("animal id from body is %v\n", card.Animal_id)
	animal, err := database.GetAnimal(pool, card.Animal_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animal", err)
		return
	}

	fmt.Printf("animal id recieved is %v\n", animal.Id)
	animalsOrganizationId, err := database.GetAnimalOrganization(pool, animal.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get organizations id from Animal", err)
		return
	}

	var notification database.Notification
	notification.Actor = userId
	notification.Notifier = animalsOrganizationId
	notification.Animal_id = animal.Id
	notification.Entity_text = "New message request"
	notification.Entity_type = 1
	notification.Status = "unseen"

	err = database.CreateNotification(pool, notification)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "issue creating notification", err)
		return
	}

	err = database.AddUserAnimalSeen(pool, userId, animal.Id, card.Liked)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not add to animals user has seen", err)
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
	return card.Animal_id
}

func ResetSeenProgress(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	userId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	err = database.ResetSeenProgress(pool, userId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not reset users card progress", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
