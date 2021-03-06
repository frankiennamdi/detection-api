package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthCheckHandler(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/api/health-check", nil)
	req := require.New(t)
	req.NoError(err)

	requestRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(StatusHandler)
	handler.ServeHTTP(requestRecorder, request)
	req.Equal(http.StatusOK, requestRecorder.Code)
	req.Equal(`{"result":"success"}`, requestRecorder.Body.String())
}
