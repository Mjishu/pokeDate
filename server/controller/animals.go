package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func AnimalController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket string, s3Client *s3.Client) {
	SetHeader(w)

	switch r.Method {
	case http.MethodPost:
		CreateAnimal(w, r, pool)
	case http.MethodPut: //* Update animal area
		UpdateAnimal(w, r, pool)
	case http.MethodDelete:
		DeleteAnimal(w, r, pool, jwtSecret, s3Bucket, s3Client)
	}
}

func ShotController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(database.GetAllShots(pool)); err != nil {
			http.Error(w, "unable to encode response", http.StatusInternalServerError)
		}
	case http.MethodDelete:
		var shot database.NewAnimalShot
		err := checkBody(w, r, &shot)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "could not find shot in body", err)
			return
		}

		err = database.DeleteAnimalShots(pool, shot)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "could not delete shot", err)
			return
		}
		respondWithJSON(w, http.StatusOK, nil)
	}
}

func GetImagePublicId(image_url string) string {
	splitString := strings.Split(image_url, "\\")
	finalString := strings.Split(splitString[len(splitString)-1], ".")

	return finalString[0]
}

func CreateAnimal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("The else in method post was called")

	var animal database.Animal
	err := checkBody(w, r, &animal)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animal in body", err)
		return
	}

	database.InsertAnimal(pool, animal)

	animal_id := database.GetAnimalByName(pool, animal.Name)

	response := map[string]interface{}{
		"Animal_id": animal_id,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func UpdateAnimal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	SetHeader(w)

	var updatedAnimal database.Animal
	err := checkBody(w, r, &updatedAnimal)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animal in body", err)
		return
	}
	fmt.Printf("updated animal is %v\n", updatedAnimal)
	for _, values := range updatedAnimal.Shots {
		newShot := database.NewAnimalShot{Animal_id: updatedAnimal.Id, Shot_id: values.Id, Date_given: values.Date_given, Date_due: values.Next_due}
		database.InsertAnimalShots(pool, newShot)
	}
	err = database.UpdateAnimal(pool, updatedAnimal)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "cannot update animal", err)
		return
	}
}

func GetAnimal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret string) {
	animalId := r.PathValue("animalID")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	_, err = auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not the correct JWT", err)
		return
	}

	animal, err := database.GetAnimal(pool, animalId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animals of this org", err)
		return
	}

	respondWithJSON(w, http.StatusOK, animal)
}

type AnimalId struct {
	Id uuid.UUID `json:"id"`
}

func DeleteAnimal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket string, s3Client *s3.Client) {
	SetHeader(w)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find JWT", err)
		return
	}

	_, err = auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	var animalId AnimalId // !changed this from checkbody might break?
	err = checkBody(w, r, &animalId)
	fmt.Printf("animal id is %v\n", animalId.Id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find id in body", err)
		return
	}

	storedAnimal, err := database.GetAnimal(pool, animalId.Id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get stored animal from database", err)
		return
	}

	err = database.DeleteAnimal(pool, animalId.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete animal", err)
		return
	}

	if storedAnimal.Image_src != nil {
		fmt.Println("image src for deleted animal is not nil")
		err = DeleteS3Object(w, r, s3Bucket, *storedAnimal.Image_src, "animals", s3Client) //! Issue with the key (image_src)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not delete image in s3", err)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func GetFrontendURL() string {
	err := godotenv.Load("../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			fmt.Printf("could not load .env%v\n", err)
		}
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatal("DATABASE_URL not set in .env")
	}
	return frontendURL
}
