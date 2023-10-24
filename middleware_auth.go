package main

import (
	"github.com/ahmadjlallad/rssagg-go/internal/auth"
	"github.com/ahmadjlallad/rssagg-go/internal/database"
	"net/http"
)

type authHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer, err := auth.GetApiKey(&r.Header)
		if err != nil {
			respondWithError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), bearer)
		if err != nil {
			respondWithError(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		handler(w, r, user)
	}
}
