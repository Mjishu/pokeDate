package controller

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func HandleUserImageUpload(w http.ResponseWriter, r *http.Request, JWTToken, s3Bucket, s3Region string, s3Client *s3.Client) {
	userIdString := r.PathValue("userID")
	userId, err := uuid.Parse(userIdString)
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
	_, err = auth.ValidateJWT(token, JWTToken) //gets user id
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not a valid JWT", err)
		return
	}

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("profile_image")
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
	if !(mimeType == "image/jpeg" || mimeType == "image/jpg" || mimeType == "image/webp" || mimeType == "image/png") { //* add more MIMETYPE here if i need
		respondWithError(w, http.StatusBadRequest, "invalid mime type", err)
		return
	}

	extensionArr := strings.Split(mediaType, "/")
	extension := extensionArr[len(extensionArr)-1]

	// create temp file here
	tempFile, err := os.CreateTemp("", "profile-pic-upload."+extension)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create temp file", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to read image data", err)
		return
	}

	//copy data
	_, err = io.Copy(tempFile, bytes.NewReader(fileData))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not copy data", err)
		return
	}

	userData, err := database.GetUserById(userId)
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
	folder_path := "profile_pictures"
	key := path.Join(folder_path, base64.RawURLEncoding.EncodeToString(newByte)+extension)

	tempFile.Seek(0, io.SeekStart)

	input := &s3.PutObjectInput{
		Bucket:      &s3Bucket,
		Key:         &key,
		Body:        tempFile,
		ContentType: &mimeType,
	}
	_, err = s3Client.PutObject(context.TODO(), input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not store image", err)
		return
	}

	imageURL := "https://" + s3Bucket + ".s3." + s3Region + ".amazonaws.com/" + key
	userData.Profile_picture = &imageURL
	err = database.UpdateUser(userData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, userData)
}
