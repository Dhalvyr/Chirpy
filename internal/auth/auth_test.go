package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "12345"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("%v\n",err)
	}

	check, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if check != true {
		t.Errorf("Hashed passord and original password do not match")
	}
}

func TestJWT(t *testing.T) {
	id := uuid.New()
	secretString := "secretString1"
	tokenString, err := MakeJWT(id, secretString, 15 * time.Minute)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	returnedID, err := ValidateJWT(tokenString, secretString)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if id != returnedID {
		t.Errorf("validated id does not match original id")
	}
}

func TestJWTtimedout(t *testing.T) {
	id := uuid.New()
	secretString := "secretString1"
	tokenString, err := MakeJWT(id, secretString, -1 * time.Second)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	_, err = ValidateJWT(tokenString, secretString)
	if err == nil {
		t.Errorf("expected expired token error, got none")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	id := uuid.New()
	secretString1 := "secretString1"
	secretString2 := "secretString2"
	tokenString, err := MakeJWT(id, secretString1, 15 * time.Minute)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	_, err = ValidateJWT(tokenString, secretString2)
	if err == nil {
		t.Errorf("expected error due to mismatching secret token, got none")
	}
}