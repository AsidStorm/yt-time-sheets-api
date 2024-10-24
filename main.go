package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"yandex.tracker.api/api"
	v1 "yandex.tracker.api/api/v1"
	"yandex.tracker.api/services/config_cache"
)

func main() {
	appCtx := &ctx{
		session: session{},
		svs: &svs{
			configCache: config_cache.Make(),
		},
	}

	server := api.NewServer(mux.NewRouter())
	v1.Register(server.Router, appCtx)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                                                                      // All origins
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions}, // Allowing only get, just an example
		AllowedHeaders: []string{"X-Auth-Token", "X-Org-ID", "X-I-Am-Token"},
	})

	log.Fatal(http.ListenAndServe(":1121", c.Handler(server.Router)))
}
