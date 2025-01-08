package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func OrganizationController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	SetHeader(w)

	switch r.URL.Path {
	case "/organizations/login":
		HandleOrganizationLogin(w, r, pool, jwtSecret)
	case "/organizations/create":
		HandleOrganizationCreate(w, r, pool)
	case "/organizations/animals":
		HandleAnimals(w, r, pool)
	}
}

func HandleAnimals(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")

		hasId, id := checkForBodyItem("id", w, r)
		animal, err := database.GetAnimal(pool, id)
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
		if err := json.NewEncoder(w).Encode(database.GetAllAnimals(pool)); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
		return
	}
}

func HandleOrganizationCreate(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	var organization database.Organization
	checkBody(w, r, &organization)

	hashedPassword, err := auth.HashPassword(organization.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not hash password", err)
		return
	}

	organization.Password = hashedPassword
	err = database.CreateOrganization(pool, organization)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create organization", err)
		return
	}

	storedOrg := database.GetOrganizationByName(pool, organization.Name)
	if (storedOrg == database.Organization{}) {
		respondWithError(w, http.StatusInternalServerError, "could not get created organization", err)
		return
	}

	respondWithJSON(w, http.StatusOK, storedOrg)
}

func HandleOrganizationLogin(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	// logic for creating login
}

func GetCurrentOrganization(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	orgId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	storedOrganization := database.GetOrganization(pool, orgId)
	if (storedOrganization == database.Organization{}) {
		respondWithError(w, http.StatusBadRequest, "no organization stored with that id", err)
		return
	}

	respondWithJSON(w, http.StatusOK, storedOrganization)
}

func GetCurrentOrganizationAnimals(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	orgId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not the correct JWT", err)
		return
	}

	animals, err := database.GetOrganizationAnimals(pool, orgId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animals of this org", err)
		return
	}

	respondWithJSON(w, http.StatusOK, animals)
}
