package services

import (
	"log"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Service is up")
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
