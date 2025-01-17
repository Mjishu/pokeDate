package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func CreateNotification(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	recipientIdString := r.PathValue("recipientID")
	recipientId, err := uuid.Parse(recipientIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find recipient ID", err)
		return
	}

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
	notification.Actor = fromId
	notification.Notifier = recipientId
	notification.Entity_text = "New message request"
	notification.Entity_type = 1
	notification.Status = "unseen"

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
