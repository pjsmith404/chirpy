package main

import (
	"github.com/pjsmith404/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
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

	err = cfg.db.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to revoke token",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
