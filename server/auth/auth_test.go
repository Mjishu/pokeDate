package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPasswordHash(t *testing.T) {
	password := "encryptMe"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf(`HashPassword("encryptMe") = %v, %v, want match for, nil\n`, hashed, err)
	}
	if err := CheckPasswordHash(password, hashed); err != nil {
		t.Fatalf(`CheckPasswordHash(password, hash) = %v, want nil`, err)
	}
}

func TestGetBearerToken(t *testing.T) { //expects http.Header but im passing string. How to to turn token into a headers
	mockHeaders := http.Header{
		"Authorization": {"Bearer ababab"},
	}
	want := "ababab"
	got, err := GetBearerToken(mockHeaders)
	if err != nil || want != got {
		t.Fatalf("GetBearerToken, want = %v, nil got %v, %v", want, got, err)
	}
}

func TestJWT(t *testing.T) {
	secret := "imASecret"
	userId, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("error creating random uuid")
		return
	}

	jwtToken, err := MakeJWT(userId, secret, 1*time.Hour)
	if err != nil {
		t.Fatalf("error creating jwt token: %v\n", err)
	}
	returnedId, err := ValidateJWT(jwtToken, secret)
	if err != nil || returnedId != userId {
		t.Fatalf("Given secret %v, want: %v, got %v\n", secret, userId, returnedId)
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "imASecret"
	userId, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("error creating random uuid")
		return
	}

	jwtToken, err := MakeJWT(userId, secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("error creating jwt token: %v\n", err)
	}
	returnedId, err := ValidateJWT(jwtToken, secret)
	if (err == nil || returnedId != uuid.UUID{}) {
		t.Fatalf("Given secret %v, want: %v, got %v\n. with error: %v\n", secret, userId, returnedId, err)
	}
}

func TestInvalidSecret(t *testing.T) {
	realSecret := "imASecret"
	fakeSecret := "hahaImFake"
	userId, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("error creating random uuid")
		return
	}

	jwtToken, err := MakeJWT(userId, realSecret, 1*time.Hour) 
	if err != nil {
		t.Fatalf("error creating jwt token: %v\n", err)
	}
	returnedId, err := ValidateJWT(jwtToken, fakeSecret) 
	if (err == nil || returnedId != uuid.UUID{}) {
		t.Fatalf("Given real secret: %v, and fake secret: %v\n want: %v, %v, got %v, nil\n", realSecret, fakeSecret, uuid.UUID{}, err, returnedId)
	}
}
