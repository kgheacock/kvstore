package router

import (
	"net/http"

	"../kvstore"

	"github.com/gorilla/mux"
)

func CreateRouter(s *kvstore.Store) *mux.Router {
	router := mux.NewRouter()

	//route registration
	router.Handle("/kv-store/{key}", wrap(s.DeleteHandler)).Methods("DELETE")
	router.Handle("/kv-store/{key}", wrap(s.PutHandler)).Methods("PUT")
	router.Handle("/kv-store/{key}", wrap(s.GetHandler)).Methods("GET")

	router.Use(loggingMiddleware)

	return router
}

func wrap(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(handler)
}
