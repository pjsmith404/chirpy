package main

import (
	"encoding/json"
	"fmt"
	"github.com/pjsmith404/chirpy/internal/auth"
	"github.com/pjsmith404/chirpy/internal/database"
	"net/http"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
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

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't get bearer token",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Failed to validate JWT",
			fmt.Errorf(`Failed to validate JWT: %v`, err),
		)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't hash password",
			err,
		)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{userID, params.Email, hash})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to update user",
			err,
		)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
