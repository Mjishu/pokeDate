package controller

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func CreateNotification(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, entity_text string, entity_type int) {
	// BODY needs animal_id and notifier

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get JWT", err)
		return
	}

	fromId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	// setup notification
	var notification database.Notification

	err = checkBody(w, r, &notification)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "did not find animal id in body", err)
		return
	}

	notification.Actor = fromId
	notification.Entity_text = entity_text
	notification.Entity_type = entity_type //1=req, 2=reply, 3=alert, 4= news
	notification.Status = "unseen"

	fmt.Printf("notification is %v\n", notification)
	err = database.CreateNotification(pool, notification)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create notification", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func GetNotifications(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get JWT", err)
		return
	}

	idToGet, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	notifications, err := database.GetNotification(pool, idToGet)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get notifications", err)
		return
	}

	respondWithJSON(w, http.StatusOK, notifications)
}

func UpdateNotifications(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get JWT", err)
		return
	}

	currentId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}
	fmt.Printf("id %v\n", currentId)
	/*
		check if notification type is request
		body should have accepted or declined when this is called
		if declined then update the status in notifications to declined
		if accepted create a new message with the actor and notifier
		send the original actor a notification saying they have new message
	*/
}
