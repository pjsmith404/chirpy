package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
	"net/http"
	"strings"
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

func TestGetBearerToken(t *testing.T) {
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

	headers := http.Header{}
	bearer := strings.Join([]string{"Bearer", jwt}, " ")
	headers.Add("Authorization", bearer)

	bearerToken, err := GetBearerToken(headers)
	if err != nil {
		t.Fatal(err)
	}

	if bearerToken != jwt {
		t.Errorf(`GetBearerToken(headers) = %q, want %q`, bearerToken, jwt)
	}
}

func TestMakeRefreshToken(t *testing.T) {
	refreshToken := MakeRefreshToken()

	if len(refreshToken) != 64 {
		t.Errorf(`MakeRefreshToken() = %v, want a 64 char string`, refreshToken)
	}
}

