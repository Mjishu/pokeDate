package controller

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func handleAnimalImageUpload(w http.ResponseWriter, r *http.Request, JWTToken string, assetsRoot string) {
	animalIdString := r.PathValue("animalID")
	animalId, err := uuid.Parse(animalIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animal id ", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find bearer token", err)
		return
	}

	//* check to make sure its an org and not a regular user
	_, err = auth.ValidateJWT(token, JWTToken) //gets org id
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not a valid JWT", err)
		return
	}

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to parse form file", err)
		return
	}
	defer file.Close()

	mediaType := header.Header.Get("Content-Type")
	mimeType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to get mimeType from header", err)
		return
	}
	if !(mimeType == "jpeg" || mimeType == "jpg" || mimeType == "webp" || mimeType == "png") { //* add more MIMETYPE here if i need
		respondWithError(w, http.StatusBadRequest, "invalid mime type", err)
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to read image data", err)
		return
	}

	animalData, err := database.GetAnimal(animalId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to get animal", err)
		return
	}

	// makes a random id
	newByte := make([]byte, 32)
	_, err = rand.Read(newByte)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to make random byte", err)
		return
	}
	randomId := base64.RawURLEncoding.EncodeToString(newByte)

	extensionArr := strings.Split(mediaType, "/")
	imagePath := randomId + "." + extensionArr[len(extensionArr)-1]
	fp := filepath.Join(assetsRoot, imagePath)

	openFile, err := os.Create(fp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create new image file", err)
		return
	}
	defer openFile.Close()

	_, err = io.Copy(openFile, bytes.NewBuffer(fileData))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not copy data to file", err)
		return
	}

	animalData.Image_src = &imagePath
	err = database.UpdateAnimal(animalData)

	respondWithJSON(w, http.StatusOK, animalData)
}
