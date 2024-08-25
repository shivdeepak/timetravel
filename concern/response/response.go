package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rainbowmga/timetravel/concern/logging"
	"github.com/rainbowmga/timetravel/model"
	"github.com/rs/zerolog/log"
)

var (
	ErrInternal = errors.New("internal error")
)

// writeJSON writes the data as json.
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	return err
}

// writeError writes the message as an error
func WriteError(w http.ResponseWriter, message string, statusCode int) error {
	log.Printf("response errored: %s", message)
	return WriteJSON(
		w,
		map[string]string{"error": message},
		statusCode,
	)
}

func WriteRecord(w http.ResponseWriter, record model.Record) {
	recordJson, err := record.ToJSON()
	if err != nil {
		err := WriteError(w, "internal error", http.StatusInternalServerError)
		logging.LogError(err)
		return

	}
	err = WriteJSON(w, recordJson, http.StatusOK)
	if err != nil {
		err := WriteError(w, "internal error", http.StatusInternalServerError)
		logging.LogError(err)
		return
	}
	logging.LogError(err)
}

func WriteRecords(w http.ResponseWriter, records []model.Record) {
	recordsJson := make([]interface{}, len(records))
	for i, record := range records {
		recordJson, err := record.ToJSON()
		if err != nil {
			err := WriteError(w, "internal error", http.StatusInternalServerError)
			logging.LogError(err)
			return
		}
		recordsJson[i] = recordJson
	}

	err := WriteJSON(w, recordsJson, http.StatusOK)
	if err != nil {
		err := WriteError(w, "internal error", http.StatusInternalServerError)
		logging.LogError(err)
		return
	}
	logging.LogError(err)
}
