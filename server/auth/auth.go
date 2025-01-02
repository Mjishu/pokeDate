package auth

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt" // go get this (should be this one https://github.com/golang-jwt/jwt
)

// go get uuid?

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
			 jwt.RegisteredClaims{
				 "iss": "pokeFind-api",
				 "iat": time.currentTime //utc
				 "exp": time.currentTime + expiresIn
				 "sub" : string(userID) //stringify user id
			 })
	tokenString, err := token.SignedString(tokenSecret)
}

// further implement the validate function
// func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID,error) {
// 	token, err := jwt.ParseWithClaims(tokenString, )
// }
