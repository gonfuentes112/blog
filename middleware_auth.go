package main

import (
	"fmt"
	"net/http"

	"github.com/gonfuentes112/blog/internal/auth"
	"github.com/gonfuentes112/blog/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *dbConfig) middlwareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.UserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
