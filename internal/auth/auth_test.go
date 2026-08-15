package auth

import "testing"

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