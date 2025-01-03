package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func RefreshToken(w http.ResponseWriter, r *http.Request, jwtSecret string) {
	SetHeader(w)

	switch r.Method {
	case http.MethodPost:
		refresh_token, err := auth.GetBearerToken(r.Header) //! this should be the refresh_token NOT jwt
		if err != nil {
			http.Error(w, "error finding refresh token", http.StatusBadRequest)
			fmt.Printf("error finding refresh_token %v\n", err)
			return
		}

		exists, userId := database.GetRefreshToken(refresh_token)
		if !exists {
			http.Error(w, "refresh token not valid", http.StatusUnauthorized)
			fmt.Printf("userId from refreshToken %v\n refreshToken exists %v\n", userId, exists)
			return
		}
		newToken, err := auth.MakeJWT(userId, jwtSecret, time.Duration(1*time.Hour))
		if err != nil {
			http.Error(w, "error creating jwt auth token", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]interface{}{"token": newToken}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

	}
}
