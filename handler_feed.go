package main

import (
	"encoding/json"
	"github.com/ahmadjlallad/rssagg-go/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type feedBody struct {
		Name string `json:"name,omitempty"`
		Url  string `json:"url,omitempty"`
	}
	body := feedBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}
	feed, err := apiCfg.DB.CreateFeed(
		r.Context(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			Name:      body.Name,
			Url:       body.Url,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
		},
	)
	if err != nil {
		respondWithError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, sendFeed(feed))
}

func (apiCfg *apiConfig) handlerListFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.ListFeeds(r.Context())
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var feedsRes []Feed
	for _, feed := range feeds {
		feedsRes = append(feedsRes, sendFeed(feed))
	}

	respondWithJson(w, http.StatusOK, feedsRes)
}
