package app

import (
	"log"
	"net/http"

	"github.com/frankiennamdi/detection-api/support"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf(support.Info, "service is up")
	responseJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
