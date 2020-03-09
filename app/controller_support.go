package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/frankiennamdi/detection-api/support"
)

func errorResponse(w http.ResponseWriter, code int, msg string) {
	responseJSON(w, code, map[string]string{"error": msg})
}

func responseJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := w.Write(response)

	if err != nil {
		log.Printf(support.Fatal, err)
	}
}
