package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gonfuentes112/blog/internal/database"

	"github.com/google/uuid"
)

func (cfg *dbConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	u1 := uuid.New()

	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	name := params.Name

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        u1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (cfg *dbConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))

}
func (cfg *dbConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get posts")
		return
	}
	respondWithJSON(w, 200, databasePostsToPosts(posts))

}
