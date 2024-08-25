package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/api"
	"github.com/rainbowmga/timetravel/logging"
	"github.com/rainbowmga/timetravel/middleware"
	"github.com/rainbowmga/timetravel/model"
	"github.com/rainbowmga/timetravel/service"
	"github.com/rs/zerolog/log"
)

// logError logs all non-nil errors
func logError(err error) {
	if err != nil {
		log.Debug().Msg("Here")
		log.Error().Err(err).Msg("")
	}
}

func main() {
	logging.InitLogging()
	model.InitDb()

	router := mux.NewRouter()

	service := service.NewSQLiteRecordService()
	api := api.NewAPI(&service)

	apiRoute := router.PathPrefix("/api/v1").Subrouter()
	apiRoute.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})
	api.CreateRoutes(apiRoute)

	loggedRouter := middleware.AccessLogMiddleware(router)

	address := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Msgf("listening on http://%s", address)
	err := srv.ListenAndServe()
	log.Fatal().Err(err).Msg("")
}
