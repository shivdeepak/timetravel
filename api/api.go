package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/api/v1"
	"github.com/rainbowmga/timetravel/concern/logging"
	"github.com/rainbowmga/timetravel/service"
)

type API struct {
	records service.RecordService
}

func NewAPI(records service.RecordService) *API {
	return &API{records}
}

// generates all api routes
func (a *API) CreateRoutes(routes *mux.Router) {
	routes.Path("/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(
				map[string]bool{"ok": true},
			)
			logging.LogError(err)
		},
	)

	apiV1 := v1.NewV1API(a.records)
	routerV1 := routes.PathPrefix("/v1").Subrouter()
	apiV1.CreateRoutes(routerV1)
}
