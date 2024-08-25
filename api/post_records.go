package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/logging"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// POST /records/{id}
// if the record exists, the record is updated.
// if the record doesn't exist, the record is created.
func (a *API) PostRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logging.LogError(err)
		return
	}

	var body map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		err := writeError(w, "invalid input; could not parse json", http.StatusBadRequest)
		logging.LogError(err)
		return
	}

	// first retrieve the record
	record, err := a.records.GetRecord(
		ctx,
		uint(idNumber),
	)

	if err == nil {
		log.Info().Msg("Update Existing Record")
		record, err = a.records.UpdateRecord(ctx, record, body)
		if err == nil {
			writeRecord(w, record)
		} else {
			err := writeError(w, ErrInternal.Error(), http.StatusInternalServerError)
			logging.LogError(err)
		}
	} else if err == gorm.ErrRecordNotFound {
		log.Info().Msg("Create New Record")
		record, err = a.records.CreateRecord(ctx, uint(idNumber), body)
		if err == nil {
			writeRecord(w, record)
		} else {
			err := writeError(w, ErrInternal.Error(), http.StatusInternalServerError)
			logging.LogError(err)
		}
	} else {
		err := writeError(w, ErrInternal.Error(), http.StatusInternalServerError)
		logging.LogError(err)
	}
}
