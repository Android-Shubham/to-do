package main

import (
	"net/http"
)

func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w,http.StatusInternalServerError, "An unexpected error occurred")
}