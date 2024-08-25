package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/concern/logging"
	"github.com/rainbowmga/timetravel/concern/response"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// POST /records/{id}
// if the record exists, the record is updated.
// if the record doesn't exist, the record is created.
func (a *API_V1) PostRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := response.WriteError(
			w,
			"invalid id; id must be a positive number",
			http.StatusBadRequest,
		)
		logging.LogError(err)
		return
	}

	var body map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		err := response.WriteError(
			w,
			"invalid input; could not parse json",
			http.StatusBadRequest,
		)
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
			response.WriteRecord(w, record)
		} else {
			err := response.WriteError(
				w,
				response.ErrInternal.Error(),
				http.StatusInternalServerError,
			)
			logging.LogError(err)
		}
	} else if err == gorm.ErrRecordNotFound {
		log.Info().Msg("Create New Record")
		record, err = a.records.CreateRecord(ctx, uint(idNumber), body)
		if err == nil {
			response.WriteRecord(w, record)
		} else {
			err := response.WriteError(
				w,
				response.ErrInternal.Error(),
				http.StatusInternalServerError,
			)
			logging.LogError(err)
		}
	} else {
		err := response.WriteError(
			w,
			response.ErrInternal.Error(),
			http.StatusInternalServerError,
		)
		logging.LogError(err)
	}
}
