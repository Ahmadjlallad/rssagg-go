package main

import (
	"net/http"
)

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, r, http.StatusBadRequest, "Something went wrong")
}
