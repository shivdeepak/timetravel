package v2

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/concern/logging"
	"github.com/rainbowmga/timetravel/concern/response"
)

// GET /records/{id}
// GetRecord retrieves the record.
func (a *API_V2) GetRecords(w http.ResponseWriter, r *http.Request) {
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

	record, err := a.records.GetRecord(
		ctx,
		uint(idNumber),
	)

	if err != nil {
		err := response.WriteError(
			w,
			fmt.Sprintf("record of id %v does not exist", idNumber),
			http.StatusBadRequest,
		)
		logging.LogError(err)
		return
	}

	response.WriteRecord(w, record)
}
