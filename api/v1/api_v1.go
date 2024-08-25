package v1

import (
	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/service"
)

type API_V1 struct {
	records service.RecordService
}

func NewV1API(records service.RecordService) *API_V1 {
	return &API_V1{records}
}

func (a *API_V1) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}
