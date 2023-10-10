package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	if code > 499 {
		log.Fatal("Responding with 5XX Error:", msg)
	}

	respondWithJson(w, code, errResponse{Error: msg})
}

type errResponse struct {
	Error string `json:"error"`
}
