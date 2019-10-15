package router

import (
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/kvstore"
	"github.com/gorilla/mux"
)

func CreateRouter(s *kvstore.Store) *mux.Router {
	router := mux.NewRouter()

	//route registration
	router.Handle("/kv-store/{key}", wrap(s.DeleteHandler)).Methods("DELETE")
	router.Use(loggingMiddleware)

	return router
}

func wrap(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(handler)
}
