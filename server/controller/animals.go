package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func GetImagePublicId(image_url string) string {
	splitString := strings.Split(image_url, "\\")
	finalString := strings.Split(splitString[len(splitString)-1], ".")

	return finalString[0]
}

func CreateAnimal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("The else in method post was called")

	animal := GetAnimalFromBody(w, r)

	database.InsertAnimal(pool, animal)

	animal_id := database.GetAnimalByName(pool, animal.Name)

	for _, values := range animal.Shots {
		newShot := database.NewAnimalShot{Animal_id: animal_id, Shot_id: values.Shot_id, Date_given: values.Date_given, Date_due: values.Date_due}
		database.InsertAnimalShots(pool, newShot)
	}

	response := map[string]interface{}{
		"Animal_id": animal_id,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func MainAnimalOperations(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	switch r.Method {
	case http.MethodPost:
		CreateAnimal(w, r, pool)
	case http.MethodPut:
		w.Header().Set("Content-Type", "application/json")

		updatedAnimal := GetUpdatedAnimalFromBody(w, r)

		for i, values := range updatedAnimal.Shots {
			fmt.Printf("Shot number: %v\n", i)
			newShot := database.NewAnimalShot{Animal_id: updatedAnimal.Id, Shot_id: values.Shot_id, Date_given: values.Date_given, Date_due: values.Date_due}
			fmt.Printf("the new shot is: %v\n", newShot)
			database.InsertAnimalShots(pool, newShot)
		}
		fmt.Printf("the animal given is %v\n", updatedAnimal)
		//!!
		database.UpdateAnimal2(pool, updatedAnimal) //TODO UPDATE THIS TO USE UPDATEANIMAL NOT 2
		//!!
		fmt.Fprintf(w, "Animal updated successfully")
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

func AnimalController(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, jwtSecret, s3Bucket string, s3Client *s3.Client) {
	SetHeader(w)

	switch r.URL.Path {
	case "/animals":
		MainAnimalOperations(w, r, pool)
	case "/animals/delete":
		switch r.Method {
		case http.MethodDelete:
			DeleteAnimal(w, r, pool, jwtSecret, s3Bucket, s3Client)
		}
	}
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

	hasId, id := checkForBodyItem("id", w, r)
	if !hasId {
		respondWithError(w, http.StatusBadRequest, "could not find id in body", err)
		return
	}

	storedAnimal, err := database.GetAnimal(pool, id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get stored animal from database", err)
		return
	}

	err = database.DeleteAnimal(pool, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete animal", err)
		return
	}

	if storedAnimal.Image_src != nil || *storedAnimal.Image_src != "" {
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

func GetAnimalFromBody(w http.ResponseWriter, r *http.Request) database.NewAnimal {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error trying to read body: %v\n", err)
		http.Error(w, "unable to read body: %v\n", http.StatusInternalServerError)
		return database.NewAnimal{}
	}
	defer r.Body.Close()

	var animal database.NewAnimal
	err = json.Unmarshal(body, &animal)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return database.NewAnimal{}
	}

	return animal
}

func GetUpdatedAnimalFromBody(w http.ResponseWriter, r *http.Request) database.UpdateAnimalStruct {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return database.UpdateAnimalStruct{}
	}
	defer r.Body.Close()

	var animal database.UpdateAnimalStruct
	err = json.Unmarshal(body, &animal)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return database.UpdateAnimalStruct{}
	}
	// fmt.Printf("Animal inside GAFB is %v\n", animal)
	return animal
}
