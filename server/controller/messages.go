package controller

import (
	"fmt"
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

func CreateConversation(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	userId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unabled to validate JWT", err)
		return
	}

	var bodyId uuid.UUID //? not sure if this should be animals id and then find the org based on the animal or the org id
	err = checkBody(w, r, bodyId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "did not find recipeint or animal id in body", err)
	}

	conversationId, err := database.CreateConversation(pool, userId, bodyId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create conversation", err)
		return
	}

	respondWithJSON(w, http.StatusOK, conversationId)
}

func CreateMessage(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
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

	userId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unabled to validate JWT", err)
		return
	}

	var messageBody database.Messages
	err = checkBody(w, r, &messageBody) //! ends up as empty string in db?
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "did not find message body", err)
	}

	messageBody.Conversation_id = messageId
	messageBody.From_id = userId

	fmt.Printf("the message body is %v\n", messageBody)
	err = database.CreateMessage(pool, userId, messageBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create message", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
