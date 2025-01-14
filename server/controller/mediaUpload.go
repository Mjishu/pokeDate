package controller

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjishu/pokeDate/auth"
	"github.com/mjishu/pokeDate/database"
)

func HandleUserImageUpload(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, JWTToken, s3Bucket, s3Region string, s3Client *s3.Client) {
	userIdString := r.PathValue("userID")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find user id ", err)
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

	//* temp file creation
	mimeType, _, tempFile, err := CreateImage(w, r, "profile_image", "profile_picture_src")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating Image", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	userData, err := database.GetUserById(pool, userId)
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
	key := path.Join(folder_path, base64.RawURLEncoding.EncodeToString(newByte))

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
	err = database.UpdateUser(pool, userData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, userData)
}

func CreateImage(w http.ResponseWriter, r *http.Request, image_key, temp_file string) (string, string, *os.File, error) {
	// reusable
	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile(image_key)
	if err != nil {
		return "", "", nil, err
	}
	defer file.Close()

	mediaType := header.Header.Get("Content-Type")
	mimeType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return "", "", nil, err
	}
	if !(mimeType == "image/jpeg" || mimeType == "image/jpg" || mimeType == "image/webp" || mimeType == "image/png") { //* add more MIMETYPE here if i need
		return "", "", nil, errors.New("mimetype does not match")
	}

	extensionArr := strings.Split(mediaType, "/")
	extension := extensionArr[len(extensionArr)-1]

	// create temp file here
	tempFile, err := os.CreateTemp("", temp_file+"."+extension)
	if err != nil {
		return "", "", nil, err
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", "", nil, err
	}

	//copy data
	_, err = io.Copy(tempFile, bytes.NewReader(fileData))
	if err != nil {
		return "", "", nil, err
	}
	return mimeType, extension, tempFile, nil
}

// ! ISSUE: doesn't properly delete the OBJECT, the objects get stored as key + extension so "a34scjpeg" instead of "a34sc.jpeg"
func DeleteS3Object(w http.ResponseWriter, r *http.Request, s3Bucket, url, prefix string, s3Client *s3.Client) error {
	fmt.Printf("url recieved is %v\n", url)
	keySplit := strings.Split(url, "/")
	key := keySplit[len(keySplit)-1]

	key = prefix + "/" + key

	fmt.Printf("the key is %v\n", key)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	}
	fmt.Printf("input is %v\n %v\n", *input.Bucket, *input.Key)

	_, err := s3Client.DeleteObject(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}
