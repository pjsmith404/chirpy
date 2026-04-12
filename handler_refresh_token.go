package main

import (
	"fmt"
	"github.com/pjsmith404/chirpy/internal/auth"
	"github.com/pjsmith404/chirpy/internal/database"
	"net/http"
	"time"
)

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't get bearer token",
			err,
		)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Refresh token not is not valid",
			err,
		)
		return
	}

	err = isTokenValid(refreshToken)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Refresh token is not valid",
			err,
		)
		return
	}

	jwt, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to generate token",
			err,
		)
		return
	}

	respondWithJson(w, http.StatusOK, RefreshTokenResponse{jwt})
}

func isTokenValid(refreshToken database.RefreshToken) error {
	if refreshToken.RevokedAt.Valid == true {
		return fmt.Errorf("Refresh token revoked at %v", refreshToken.RevokedAt.Time)
	}

	if refreshToken.ExpiresAt.Valid == true && time.Now().After(refreshToken.ExpiresAt.Time) {
		return fmt.Errorf("Refresh token expired at %v", refreshToken.ExpiresAt.Time)
	}

	return nil
}
