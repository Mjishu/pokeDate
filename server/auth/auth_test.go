package auth

import (
	"testing"
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
	headers := { // turn into http.Headers
		"Authorization": "Bearer ababab"
	}
	want := "ababab"
	got,err := GetBearerToken(token)
	if err != nil {
		t.Fatalf("GetBearerToken, want = %v, nil got %v, %v",want, got, err)
	}
}
