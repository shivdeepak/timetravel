package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/api/v1"
	"github.com/rainbowmga/timetravel/api/v2"
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
	apiV1 := v1.NewV1API(a.records)
	routerV1 := routes.PathPrefix("/v1").Subrouter()

	routerV1.Path("/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(
				map[string]bool{"ok": true},
			)
			logging.LogError(err)
		},
	)

	apiV1.CreateRoutes(routerV1)

	apiV2 := v2.NewV2API(a.records)
	routerV2 := routes.PathPrefix("/v2").Subrouter()
	apiV2.CreateRoutes(routerV2)
}
