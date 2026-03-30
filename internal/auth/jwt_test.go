package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "1234567890"
	expiresIn, err := time.ParseDuration("10s")
	if err != nil {
		t.Fatal(err)
	}

	jwt, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatal(err)
	}

	if jwt == "" {
		t.Errorf(`MakeJWT(userID, tokenSecret, ExpiresIn) = %q, want a jwt`, jwt)
	}

	validatedUserID, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Fatal(err)
	}

	if userID != validatedUserID {
		t.Errorf(`ValidateJWT(tokenString, tokenSecret) = %q, want %q`, validatedUserID, userID)
	}
}

