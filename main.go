package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gonfuentes112/blog/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	DB *database.Queries
}

func main() {
	const filepathRoot = "."
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONN")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database")
	}
	dbQueries := database.New(db)
	dbCfg := dbConfig{
		DB: dbQueries,
	}

	go startScraping(dbQueries, 10, time.Minute)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerError)
	mux.HandleFunc("POST /v1/users", dbCfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", dbCfg.middlwareAuth(dbCfg.handlerGetUser))
	mux.HandleFunc("POST /v1/feeds", dbCfg.middlwareAuth(dbCfg.handlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", dbCfg.handlerGetFeeds)
	mux.HandleFunc("POST /v1/feed_follows", dbCfg.middlwareAuth(dbCfg.handlerCreateFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", dbCfg.middlwareAuth(dbCfg.handlerGetFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", dbCfg.middlwareAuth(dbCfg.handlerDeleteFeedFollow))
	mux.HandleFunc("GET /v1/posts", dbCfg.middlwareAuth(dbCfg.handlerGetPostsForUser))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
