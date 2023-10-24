package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ahmadjlallad/rssagg-go/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	loadEnvError := loadEnv()
	if loadEnvError != nil {
		log.Fatalln(loadEnvError)
	}

	port := getEnvOrFatal("PORT")
	dbUrl := getEnvOrFatal("DB_URL")

	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	queries := database.New(connection)
	apiCfg := apiConfig{queries}

	r := chi.NewRouter()
	v1 := chi.NewRouter()
	r.Use(makeCors(), middleware.Logger)

	// v1 mounted
	r.Mount("/v1", v1)
	v1.Get("/healthz", handlerReadiness)
	r.Get("/pingHandler", pingHandler)
	v1.Get("/error", handlerErr)

	v1.Post("/users", apiCfg.handlerCreateUser)
	v1.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))

	v1.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1.Get("/feeds", apiCfg.handlerListFeeds)

	v1.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerListFeedsFollows))
	v1.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	server := &http.Server{Handler: r, Addr: ":" + port}
	fmt.Printf("server is running on http://localhost:%v", port)
	if server.ListenAndServe() != nil {
		log.Fatalln(err)
	}
}

func makeCors() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

func loadEnv() error {

	err := godotenv.Load(".env")
	if err != nil {
		return errors.New("FATAL: .env file dose not exist")
	}
	return nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Bong"))
	if err != nil {
		log.Fatalf("FATAL: Couldn't server the request %v", r.URL)
		return
	}
}

func getEnvOrFatal(env string) string {
	envValue := os.Getenv(env)
	if envValue == "" {
		log.Fatalf("FATAL: %v is not set\n", env)
	}

	return envValue
}
