package v2

import (
	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/service"
)

type API_V2 struct {
	records service.RecordService
}

func NewV2API(records service.RecordService) *API_V2 {
	return &API_V2{records}
}

func (a *API_V2) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}
