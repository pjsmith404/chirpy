package main

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

