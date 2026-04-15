package main

import (
	"database/sql"
	"encoding/json"
	"github.com/pjsmith404/chirpy/internal/auth"
	"github.com/pjsmith404/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Couldn't find user",
			err,
		)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to validate password",
			err,
		)
		return
	}

	if match != true {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid password",
			err,
		)
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to generate token",
			err,
		)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	now := time.Now()
	refreshTokenExpiry := now.Add(time.Hour * 24 * 60)
	refreshTokenParams := database.CreateRefreshTokenParams{
		refreshToken,
		user.ID,
		sql.NullTime{refreshTokenExpiry, true},
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to create refresh token",
			err,
		)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			user.ID,
			user.CreatedAt,
			user.UpdatedAt,
			user.Email,
			user.IsChirpyRed,
		},
		Token:        jwt,
		RefreshToken: refreshToken,
	})
}
