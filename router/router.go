package router

import (
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/kvstore"

	"github.com/gorilla/mux"
)

func CreateRouter(s *kvstore.Store, h *hasher.Store) *mux.Router {
	router := mux.NewRouter()
	storeRouter := router.PathPrefix("/kv-store/keys").Subrouter()

	//route registration
	storeRouter.Handle("/{key}", wrap(s.DeleteHandler)).Methods("DELETE")
	storeRouter.Handle("/{key}", wrap(s.PutHandler)).Methods("PUT")
	storeRouter.Handle("/{key}", wrap(s.GetHandler)).Methods("GET")

	router.Handle("/kv-store/key-count", wrap(s.GetKeyCountHandler)).Methods("GET")
	router.Handle("/kv-store/view-change", wrap(s.ExternalReshardHandler)).Methods("PUT")
	router.Handle("/internal/vc-complete", wrap(s.ReshardCompleteHandler)).Methods("GET")
	router.Handle("/internal/view-change", wrap(s.InternalReshardHandler)).Methods("PUT")

	middlewareStore := NewStore(h, s)
	router.Use(middlewareStore.loggingMiddleware)
	router.Use(middlewareStore.checkSourceMiddleware)
	router.Use(middlewareStore.bufferRequestMiddleware)
	storeRouter.Use(middlewareStore.validateParametersMiddleware)
	storeRouter.Use(middlewareStore.forwardMiddleware)

	return router
}

func wrap(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(handler)
}
