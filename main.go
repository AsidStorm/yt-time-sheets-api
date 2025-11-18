package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"yandex.tracker.api/api"
	v1 "yandex.tracker.api/api/v1"
	v3 "yandex.tracker.api/api/v3"
	"yandex.tracker.api/services/config_cache"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("unable to start [backend] on logger due [%s]", err))
	}
	defer logger.Sync()

	baseCtx := &ctx{
		logger: logger.Sugar(),
		svs: &svs{
			configCache: config_cache.Make(),
		},
	}

	appCtx := baseCtx.WithSession(session{})

	server := api.NewServer(mux.NewRouter())
	v1.Register(server.Router, appCtx)
	v3.Register(server.Router, appCtx)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"X-Auth-Token", "X-Org-ID", "X-I-Am-Token"},
	})

	log.Fatal(http.ListenAndServe(":1121", c.Handler(server.Router)))
}
