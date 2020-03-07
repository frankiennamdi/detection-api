package app

import (
	"encoding/json"
	"github.com/frankiennamdi/detection-api/core"
	"log"
	"net/http"

	"github.com/frankiennamdi/detection-api/models"
	"github.com/frankiennamdi/detection-api/support"
)

//
type EventDetectionController struct {
	DetectionService core.DetectionService
}

func (controller *EventDetectionController) EventDetectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "POST Required")
	}

	var event *models.Event

	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		log.Printf(support.Error, err)
		respondWithError(w, http.StatusBadRequest, "can pass request body")

		return
	}

	suspiciousTravelResult, err := controller.DetectionService.ProcessEvent(event)
	if err != nil {
		log.Printf(support.Error, err)
		respondWithError(w, http.StatusBadRequest, "Unable to process request")

		return
	}

	respondWithJSON(w, http.StatusOK, suspiciousTravelResult)
}
