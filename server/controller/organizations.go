package controller

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func OrganizationController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket, s3Region string) {
	SetHeader(w)

	switch r.URL.Path {
	case "/organizations/animals":
		switch r.Method {
		case http.MethodPost:
			GetCurrentOrganizationAnimals(w, r, pool, jwtSecret)
		}
	case "/organizations/animals/create":
		switch r.Method {
		case http.MethodPost:
			CreateOrganizationAnimal(w, r, pool, jwtSecret)
		}
	}
}

func CreateOrganizationAnimal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
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

	fmt.Printf("animal shots are %v\n", incomingAnimal.Shots)

	for _, values := range incomingAnimal.Shots {
		fmt.Printf("adding this shot %v\n", values.Name)
		newShot := database.NewAnimalShot{Animal_id: storedAnimalId, Shot_id: values.Id, Date_given: values.Date_given, Date_due: values.Next_due}
		err = database.InsertAnimalShots(pool, newShot)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not create shot", err)
			return
		}
	}

	response := map[string]interface{}{
		"Animal_id": storedAnimalId,
	}
	respondWithJSON(w, http.StatusOK, response)

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
