package v3

import (
	"net/http"

	"yandex.tracker.api/domain"

	"github.com/gorilla/mux"
)

const RouterPrefix = "/api/v3"

func Register(r *mux.Router, c domain.Context) {
	router := r.PathPrefix(RouterPrefix).Subrouter()

	handle := wrapHandler(c, router)

	handle("/result", http.MethodPost, Result, false)
}
