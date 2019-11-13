package router

import (
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/kvstore"

	"github.com/gorilla/mux"
)

func CreateRouter(s *kvstore.Store, h *hasher.Store) *mux.Router {
	router := mux.NewRouter()

	//route registration
	router.Handle("/kv-store/{key}", wrap(s.DeleteHandler)).Methods("DELETE")
	router.Handle("/kv-store/{key}", wrap(s.PutHandler)).Methods("PUT")
	router.Handle("/kv-store/{key}", wrap(s.GetHandler)).Methods("GET")

	middlewareStore := NewStore(h, s)
	router.Use(middlewareStore.loggingMiddleware)
	router.Use(middlewareStore.validateParametersMiddleware)
	router.Use(middlewareStore.forwardMiddleware)

	return router
}

func wrap(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(handler)
}
