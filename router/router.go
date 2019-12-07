package router

import (
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/kvstore"

	"github.com/gorilla/mux"
)

func CreateRouter(s *kvstore.Store, h *hasher.Store) *mux.Router {
	router := mux.NewRouter()
	nonContextualRouter := router.PathPrefix("/kv-store").Subrouter()
	storeRouter := nonContextualRouter.PathPrefix("/keys").Subrouter()
	//route registration
	// storeRouter.Handle("/{key}", wrap(s.DeleteHandler)).Methods("DELETE")
	storeRouter.Handle("/{key}", wrap(s.PutHandler)).Methods("PUT")
	storeRouter.Handle("/{key}", wrap(s.GetHandler)).Methods("GET")

	nonContextualRouter.Handle("/key-count", wrap(s.GetKeyCountHandler)).Methods("GET")
	nonContextualRouter.Handle("/view-change", wrap(s.ExternalReshardHandler)).Methods("PUT")

	nonContextualRouter.Handle("/kv-store/shards", wrap(s.GetShardHandler)).Methods("GET")
	nonContextualRouter.Handle("/kv-store/shards/{id}", wrap(s.GetShardByIdHandler)).Methods("GET")

	router.Handle("/internal/vc-complete", wrap(s.ReshardCompleteHandler)).Methods("GET")
	router.Handle("/internal/view-change", wrap(s.InternalReshardHandler)).Methods("PUT")
	router.Handle("/internal/prepare-for-vc", wrap(s.PrepareReshardHandler)).Methods("PUT")
	router.Handle("/internal/reshard-put/{key}", wrap(s.ReshardPutHandler)).Methods("PUT")
	router.Handle("/internal/gossip-put/{key}", wrap(s.GossipPutHandler)).Methods("PUT")

	middlewareStore := NewStore(h, s)

	router.Use(middlewareStore.loggingMiddleware)
	router.Use(middlewareStore.checkSourceMiddleware)
	router.Use(middlewareStore.bufferRequestMiddleware)

	nonContextualRouter.Use(middlewareStore.passthroughCausalContext)

	storeRouter.Use(middlewareStore.checkVectorClock)
	storeRouter.Use(middlewareStore.validateParametersMiddleware)
	storeRouter.Use(middlewareStore.forwardMiddleware)

	return router
}

func wrap(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(handler)
}
