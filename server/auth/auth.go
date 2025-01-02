package auth

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt" // go get this (should be this one https://github.com/golang-jwt/jwt
)

// go get uuid?

type CustomClaims struct {
	Id uuid.UUID`json:"Id"`
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
// ADD UNIT TESTS !!!!!
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID,error) {
	token, err := jwt.ParseWithClaims(tokenString,&CustomClaims{}, func(token *jwt.Token)(interface{}, error) {
		return []byte("allyourbase), nil //default return given by docs
	}
	if err != nil {
		log.Fatal(err)
		} else if claims, ok := token.Claims.(*CustomClaims); ok { //this is my token.Claims
			fmt.Println(claims.Id, claims.RegisteredClaims.Issuer)// gets issuer be want id?	
	} else {
		log.Fatal("unknown claims type, cannot proceed")
		}		      
}

func GetBearerToken(headers http.Header) (string, error) {
	// look for auth header in headers and return the token_string from bearer
	// get the string and then split it based on spaces and then return length -1?
	if !headers.Authorization {
		return "", new error("Could not find authorization header")
		}
	bearerToken := strings.Split(headers.Authorization, " ")
	return strings.TrimSpace(bearerToken[len(bearerToken) - 1]), nil
}
			      
