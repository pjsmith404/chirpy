package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"fmt"
	"net/http"
	"strings"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error)  {
	now := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	parseCallback := func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, parseCallback)

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok != true {
		return uuid.Nil, fmt.Errorf("Failed to get claim")
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil	
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")

	if bearerToken == "" {
		return "", fmt.Errorf("No token found in Authorization header")
	}

	bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")

	return bearerToken, nil
}

