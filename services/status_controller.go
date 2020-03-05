package services

import (
	"github.com/frankiennamdi/detection-api/support"
	"log"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(support.Info, "service is up")
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
