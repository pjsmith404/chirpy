package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pjsmith404/chirpy/internal/auth"
	"github.com/pjsmith404/chirpy/internal/database"
	"net/http"
	"time"
	"fmt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type LoggedInUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create user",
			err,
		)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{params.Email, hash})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create user",
			err,
		)
		return
	}

	respondWithJson(w, http.StatusCreated, User(user))
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
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

	expiryDuration := "1h"
	secondsInHour := 3600
	if params.ExpiresInSeconds == 0 && params.ExpiresInSeconds < secondsInHour {
		expiryDuration = fmt.Sprintf(`%vs`, params.ExpiresInSeconds)
	}
	expiresIn, err := time.ParseDuration(expiryDuration)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to parse expires_in_seconds",
			err,
		)
		return
	}
	
	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to generate token",
			err,
		)
		return
	}

	respondWithJson(w, http.StatusOK, LoggedInUser{user.ID, user.CreatedAt, user.UpdatedAt, user.Email, jwt})
}
