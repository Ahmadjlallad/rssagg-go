package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers *http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("unauthorized")
	}

	bearer := strings.Split(auth, " ")
	if len(bearer) != 2 {
		return "", errors.New("unauthorized")
	}

	if bearer[0] != "Bearer" {
		return "", errors.New("unauthorized")
	}

	return bearer[1], nil
}
