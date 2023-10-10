package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	port, loadEnvError := loadEnv()
	if loadEnvError != nil {
		log.Fatalln(loadEnvError)
		return
	}

	r := chi.NewRouter()
	v1 := chi.NewRouter()
	r.Use(makeCors(), middleware.Logger)

	// v1 mounted
	r.Mount("/v1", v1)
	v1.Get("/healthz", handlerReadiness)
	v1.Get("/error", handlerErr)

	r.Get("/ping", ping)

	runServer(port, r)
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

func loadEnv() (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		return "", errors.New("FATAL: .env file dose not exist")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return "", errors.New("FATAL: port is not set")
	}

	return port, nil
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Bong"))
	if err != nil {
		log.Fatalf("FATAL: Couldn't server the request %v", r.URL)
		return
	}
}

func runServer(port string, r http.Handler) *http.Server {
	server := &http.Server{Handler: r, Addr: ":" + port}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	log.Printf("server started on %v\n", port)
	return server
}
