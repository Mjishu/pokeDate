package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/database"
)

func ShotController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	if r.URL.Path != "/shots" {
		fmt.Fprintf(w, "Incorrect Path")
		return
	}
	SetHeader(w)

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(database.GetAllShots(pool)); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
	}
}
