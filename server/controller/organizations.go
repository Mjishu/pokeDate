package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func OrganizationController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket, s3Region string) {
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
			HandleOrganizationCreate(w, r, pool, jwtSecret, s3Bucket, s3Region)
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
		switch r.Method {
		case http.MethodPost:
			GetCurrentOrganizationAnimals(w, r, pool, jwtSecret)
		}
	case "/organizations/animals/create":
		switch r.Method {
		case http.MethodPost:
			CreateNewAnimal(w, r, pool, jwtSecret)
		}
	}
}

func CreateNewAnimal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	SetHeader(w)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT in header", err)
		return
	}

	orgId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	var incomingAnimal database.Animal
	err = checkBody(w, r, &incomingAnimal)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "did not find animal in body", err)
		return
	}

	storedAnimalId, err := database.CreateNewAnimal(pool, incomingAnimal)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating new animal", err)
		return
	}

	_, err = database.CreateOrganizationAnimal(pool, orgId, storedAnimalId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create organization animal link", err)
		return
	}

	fmt.Fprintf(w, "Body does not have an id!")

}

func HandleOrganizationCreate(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket, s3Region string) {

	var organization database.Organization
	checkBody(w, r, &organization)

	hashedPassword, err := auth.HashPassword(organization.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not hash password", err)
		return
	}

	defaultPfp := "https://" + s3Bucket + ".s3." + s3Region + ".amazonaws.com/profile_pictures/default.webp"
	organization.Profile_picture = defaultPfp

	organization.Password = hashedPassword
	storedId, err := database.CreateOrganization(pool, organization) //* nil error here?
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create organization", err)
		return
	}

	//* creates jwt
	expiresIn := time.Duration(15 * time.Minute)

	token, err := auth.MakeJWT(storedId, jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create JWT", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not create refresh token", err)
		return
	}

	_, err = database.CreateRefreshToken(pool, refresh_token, storedId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not store refresh token", err)
		return
	}

	response := map[string]interface{}{
		"name":          organization.Name,
		"id":            organization.Id,
		"status":        http.StatusOK,
		"token":         token,
		"refresh_token": refresh_token,
	}

	respondWithJSON(w, http.StatusOK, response)
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
