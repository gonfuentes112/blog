package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gonfuentes112/blog/internal/database"

	"github.com/google/uuid"
)

func (cfg *dbConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	u1 := uuid.New()

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	name := params.Name

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        u1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	respondWithJSON(w, 201, databaseFeedToFeed((feed)))
}

func (cfg *dbConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}
	respondWithJSON(w, 200, databaseFeedsToFeeds((feeds)))
}
