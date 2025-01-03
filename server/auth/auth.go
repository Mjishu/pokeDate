package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// go get uuid?

type CustomClaims struct {
	Id uuid.UUID `json:"sub"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
func CheckPasswordHash(password, hash string) error {
	bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "pokeFind-api",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		fmt.Printf("there was an error trying to sign the token: %v", err)
		return "", err
	}
	return tokenString, nil
}

// further implement the validate function
// ADD UNIT TESTS !!!!!
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil //default return given by docs
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return uuid.UUID{}, errors.New("invalid token")
	}

	return claims.Id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authToken := headers.Get("Authorization") // if field has multiple: headers["Authorization"] -> returns a struct of strings
	if authToken == "" {
		return "", errors.New("could not find authorization header")
	}
	bearerToken := strings.Split(authToken, " ")
	return strings.TrimSpace(bearerToken[len(bearerToken)-1]), nil
}

func UserValid(header http.Header, jwtSecret string) error {
	bearerToken, err := GetBearerToken(header)
	if err != nil {
		return err
	}
	_, err = ValidateJWT(bearerToken, jwtSecret)
	if err != nil {
		return err
	}
	return nil
}
