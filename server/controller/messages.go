package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func CurrentUserMessages(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, JWTSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	userId, err := auth.ValidateJWT(token, JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	messages, err := database.GetConversations(pool, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get user messages", err)
		return
	}

	respondWithJSON(w, http.StatusOK, messages)
}

func GetMessage(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, JWTSecret string) {
	messageIdString := r.PathValue("messageID")
	messageId, err := uuid.Parse(messageIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not parse messageId to UUID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	_, err = auth.ValidateJWT(token, JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	message, err := database.GetConversation(pool, messageId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get message", err)
		return
	}

	//todo check : for each users in message.users if one of the users has userId
	// if message.User_id != userId {
	// 	respondWithError(w, http.StatusUnauthorized, "not the correct user", errors.New("invalid user"))
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, message)
}
