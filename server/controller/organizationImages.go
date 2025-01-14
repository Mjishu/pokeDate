package controller

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

//* check if animalImage exists first and update that image? not sure how we would do that though, would probably have to delete prev image and
//* just make a new one

func UploadAnimalImage(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, JWTToken, s3Bucket, s3Region string, s3Client *s3.Client) {
	fmt.Printf("The url path is %v\n", r.URL.Path)
	animalIdString := r.PathValue("animalID")
	animalId, err := uuid.Parse(animalIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find animalId", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find bearer token", err)
		return
	}

	_, err = auth.ValidateJWT(token, JWTToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not validate JWT", err)
		return
	}

	mimeType, extension, tempFile, err := CreateImage(w, r, "animal_image", "animal_image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating Image", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	newByte := make([]byte, 32)
	_, err = rand.Read(newByte)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to make random byte", err)
		return
	}

	folder_path := "animals"
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

	var imagePriority int

	// err = checkBody(w, r, imagePriority)
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "did not find image priority in body", err)
	// 	return
	// }

	imageURL := "https://" + s3Bucket + ".s3." + s3Region + ".amazonaws.com/" + key
	imagePriority = 0
	err = database.AddAnimalImage(pool, imageURL, animalId, imagePriority) // todo: send priority over from frontend in the req body default should be 0
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}
