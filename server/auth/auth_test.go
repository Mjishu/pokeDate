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
