package v1

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"yandex.tracker.api/domain"
	"yandex.tracker.api/presenters/jsonapi"
)

type handleFunc func(path, method string, callable handlerFunc, public bool)
type handlerFunc func(v map[string]string, c domain.Context, w http.ResponseWriter, r *http.Request) (int, []byte, error)

func wrapHandler(c domain.Context, r *mux.Router) handleFunc {
	internalWrapper := func(callable handlerFunc, public bool) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)

			ctx := c.WithSession(session{
				authToken:      r.Header.Get("X-Auth-Token"),
				organizationId: r.Header.Get("X-Org-ID"),
				iAmToken:       r.Header.Get("X-I-Am-Token"),
			})

			if !ctx.Session().IsAuthorized() && !public {
				log.Printf("%s %s: %s\n", r.Method, r.RequestURI, errors.New("not authorized"))

				JsonError(w, http.StatusUnauthorized, errors.New("not authorized"))
			} else {
				status, response, err := callable(vars, ctx, w, r)

				if err != nil {
					log.Printf("%s %s: %s\n", r.Method, r.RequestURI, err)

					JsonError(w, status, err)
				} else {
					if status == http.StatusCreated {
						JsonSuccessCreated(w, response)
					} else {
						if status == http.StatusAccepted {
							JsonSuccessAccepted(w, response)
						} else {
							JsonSuccess(w, response)
						}
					}
				}
			}
		}
	}

	return func(path, method string, callable handlerFunc, public bool) {
		r.HandleFunc(path, internalWrapper(callable, public)).Methods(method)
	}
}

func JsonError(w http.ResponseWriter, status int, err error) {
	j, _ := json.Marshal(jsonapi.NewErrorResponse(status, err))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func JsonSuccessCreated(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

func JsonSuccess(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func JsonSuccessAccepted(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(bytes)
}

func BadRequest(err error) (int, []byte, error) {
	return http.StatusBadRequest, nil, err
}

func InternalServerError(err error) (int, []byte, error) {
	return http.StatusInternalServerError, nil, err
}

func NotFound(err error) (int, []byte, error) {
	return http.StatusNotFound, nil, err
}

func OK(body []byte) (int, []byte, error) {
	return http.StatusOK, body, nil
}

func Created(body []byte) (int, []byte, error) {
	return http.StatusCreated, body, nil
}
