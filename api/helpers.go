package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rainbowmga/timetravel/logging"
	"github.com/rs/zerolog/log"
)

var (
	ErrInternal = errors.New("internal error")
)

// writeJSON writes the data as json.
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err
}

// writeError writes the message as an error
func writeError(w http.ResponseWriter, message string, statusCode int) error {
	log.Printf("response errored: %s", message)
	return writeJSON(
		w,
		map[string]string{"error": message},
		statusCode,
	)
}
