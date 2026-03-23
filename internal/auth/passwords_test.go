package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "pa$$word"
	hash, err := HashPassword(password)

	if err != nil {
		t.Fatal(err)
	}

	if hash == "" {
		t.Errorf(`HashPassword(password) = %q, want a hash`, hash)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "pa$$word"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	want := true
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatal(err)
	}

	if match != want {
		t.Errorf(`CheckPasswordHash(%v, %v) = %v, want %v`, password, hash, match, want)
	}

	wrongPassword := "12345678"
	want = false
	match, err = CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatal(err)
	}

	if match != want {
		t.Errorf(`CheckPasswordHash(%v, %v) = %v, want %v`, wrongPassword, hash, match, want)
	}
}
