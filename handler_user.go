package main

import (
	"encoding/json"
	"fmt"
	"github.com/ahmadjlallad/rssagg-go/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	bodyParams := body{}
	err := decoder.Decode(&bodyParams)
	if err != nil {
		respondWithError(w, r, 400, fmt.Sprintf("Error parsing json %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      bodyParams.Name,
	})
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, fmt.Sprintf("Couldn't create user %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, sendUser(user))
}

func (apiCfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, http.StatusOK, sendUser(user))
}
