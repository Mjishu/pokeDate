package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
	"github.com/mjishu/pokeDate/database"
)

func GetImagePublicId(image_url string) string {
	splitString := strings.Split(image_url, "\\")
	finalString := strings.Split(splitString[len(splitString)-1], ".")

	return finalString[0]
}

func MainAnimalOperations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		fmt.Println("The else in method post was called")

		animal := GetAnimalFromBody(w, r)

		// if animal.Image_src != "" {
		// 	animal_public_id := GetImagePublicId(animal.Image_src)
		// 	fmt.Printf("iamge public id %v\n and animal url is %v\n", animal_public_id, animal.Image_src)
		// 	database.UploadImage(cld, ctx, animal.Image_src)
		// 	return //! get rid of this to make everything else work
		// }

		database.InsertAnimal(animal)

		animal_id := database.GetAnimalByName(animal.Name)

		for _, values := range animal.Shots {

			newShot := database.NewAnimalShot{Animal_id: animal_id, Shot_id: values.Shot_id, Date_given: values.Date_given, Date_due: values.Date_due}
			database.InsertAnimalShots(newShot)
		}

		fmt.Fprintf(w, "Animal created Successfully!: %v\n", animal)
	case http.MethodPut:
		w.Header().Set("Content-Type", "application/json")

		updatedAnimal := GetUpdatedAnimalFromBody(w, r)

		for i, values := range updatedAnimal.Shots {
			fmt.Printf("Shot number: %v\n", i)
			newShot := database.NewAnimalShot{Animal_id: updatedAnimal.Id, Shot_id: values.Shot_id, Date_given: values.Date_given, Date_due: values.Date_due}
			fmt.Printf("the new shot is: %v\n", newShot)
			database.InsertAnimalShots(newShot)
		}
		fmt.Printf("the animal given is %v\n", updatedAnimal)
		database.UpdateAnimal(updatedAnimal)
		fmt.Fprintf(w, "Animal updated successfully")
	case http.MethodDelete:
		w.Header().Set("Content-Type", "application/json")

		hasId, id := checkForBodyItem("id", w, r)

		if hasId {
			database.DeleteAnimal(id)
			fmt.Fprintf(w, "Animal was removed successfully")
			return
		}
		fmt.Fprintf(w, "Id was not found in body!")
	}
}

func AnimalController(w http.ResponseWriter, r *http.Request, cld *cloudinary.Cloudinary, ctx context.Context) {
	SetHeader(w)

	if r.URL.Path == "/animals" {
		MainAnimalOperations(w, r)
	} else if r.URL.Path == "/animals/images" {
		AnimalImageOperations(w, r, cld, ctx)
	}
}

func AnimalImageOperations(w http.ResponseWriter, r *http.Request, cld *cloudinary.Cloudinary, ctx context.Context) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "unable to parse form data", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("Image_src")
		if err != nil {
			http.Error(w, "error trying to form file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filePath := fmt.Sprintf("./uploads/%s", handler.Filename)
		if filePath == "" {
			return
		}

		// _, image_data := checkForBodyItem("FormData", w, r)
		// database.UploadImage(cld, ctx, image_data)
	}
}

func GetFrontendURL() string {
	err := godotenv.Load("../.env")
	if err != nil {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal("ERror loading .env file")
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

func checkForBodyItem(key string, w http.ResponseWriter, r *http.Request) (bool, any) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return false, ""
	}
	defer r.Body.Close()

	var elements map[string]interface{}
	err = json.Unmarshal(body, &elements)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusInternalServerError)
		return false, ""
	}

	if _, exists := elements[key]; exists {
		return true, elements[key]
	}
	return false, ""
}
