package main

import (
	"encoding/json"
	"github.com/ahmadjlallad/rssagg-go/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type feedFollowsBody struct {
		FeedID uuid.UUID `json:"feed_id,omitempty"`
	}
	body := feedFollowsBody{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	feedFollows, err := apiCfg.DB.CreateFeedFollows(
		r.Context(),
		database.CreateFeedFollowsParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    body.FeedID,
		},
	)

	if err != nil {
		respondWithError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondWithJson(w, http.StatusCreated, sendFeedFollow(feedFollows))
}

func (apiCfg *apiConfig) handlerListFeedsFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	listFeedFollows, err := apiCfg.DB.ListFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var feedFollows []FeedFollows
	for _, feedFollow := range listFeedFollows {
		feedFollows = append(feedFollows, sendFeedFollow(feedFollow))
	}

	respondWithJson(w, http.StatusOK, feedFollows)
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowID")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{ID: feedFollowId, UserID: user.ID})
	if err != nil {
		respondWithError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJson(w, http.StatusNoContent, struct{}{})
}
