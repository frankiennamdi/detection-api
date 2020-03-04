package services

import (
	"encoding/json"
	"github.com/frankiennamdi/detection-api/models"
	"github.com/frankiennamdi/detection-api/support"
	"log"
	"net/http"
)

type EventDetectionController struct {
	DetectionService DetectionService
	EventRepository  EventRepository
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

	if _, err := controller.EventRepository.Insert([]*models.Event{event}); err != nil {
		log.Printf(support.Error, err)
		respondWithError(w, http.StatusBadRequest, "Unable to process request")

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
