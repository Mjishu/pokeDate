package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func OrganizationController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	SetHeader(w)

	switch r.URL.Path {
	case "/organizations/login":
		switch r.Method {
		case http.MethodPost:
			HandleOrganizationLogin(w, r, pool, jwtSecret)
		default:
			respondWithError(w, http.StatusBadRequest, "no route found!", nil)
		}
	case "/organizations/create":
		switch r.Method {
		case http.MethodPost:
			HandleOrganizationCreate(w, r, pool, jwtSecret)
		}
	case "/organizations/current":
		switch r.Method {
		case http.MethodPost:
			GetCurrentOrganization(w, r, pool, jwtSecret)
		}
	case "/organizations/update":
		switch r.Method {
		case http.MethodPut:
			UpdateOrganization(w, r, pool, jwtSecret)
		}
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

func HandleOrganizationCreate(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {

	//TODO check if org with name already exists

	var organization database.Organization
	checkBody(w, r, &organization)

	hashedPassword, err := auth.HashPassword(organization.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not hash password", err)
		return
	}

	organization.Password = hashedPassword
	err = database.CreateOrganization(pool, organization) //* nil error here?
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create organization", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func HandleOrganizationLogin(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	// logic for creating login
	var organization database.Organization
	checkBody(w, r, &organization)

	storedOrg := database.GetOrganizationByName(pool, organization.Name)

	err := auth.CheckPasswordHash(organization.Password, storedOrg.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "passwords do not match", err)
		return
	}

	expiresIn := time.Duration(15 * time.Minute)

	token, err := auth.MakeJWT(storedOrg.Id, jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create JWT", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create refresh token", err)
		return
	}

	_, err = database.CreateRefreshToken(pool, refresh_token, storedOrg.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not store refresh token", err)
		return
	}

	response := map[string]interface{}{
		"name":          storedOrg.Name,
		"id":            storedOrg.Id,
		"status":        http.StatusOK,
		"token":         token,
		"refresh_token": refresh_token,
	}

	respondWithJSON(w, http.StatusOK, response)
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

func UpdateOrganization(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get bearer token", err)
		return
	}

	orgId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	var incomingOrg database.Organization
	err = checkBody(w, r, &incomingOrg)
	if err != nil || incomingOrg.Id != orgId {
		respondWithError(w, http.StatusBadRequest, "error trying to get incoming org", err)
		return
	}

	storedOrg := database.GetOrganization(pool, orgId)
	if (storedOrg == database.Organization{}) {
		respondWithError(w, http.StatusBadRequest, "could not find stored orgnaization with that id", err)
		return
	}

	if storedOrg.Name != incomingOrg.Name {
		storedOrg.Name = incomingOrg.Name
	}
	if storedOrg.Email != incomingOrg.Email {
		storedOrg.Email = incomingOrg.Email
	}

	err = database.UpdateOrganization(pool, storedOrg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, storedOrg)
}
