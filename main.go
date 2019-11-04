package main

import (
	"log"
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/kvstore"
	"github.com/colbyleiske/cse138_assignment2/router"
)

func main() {
	config.GenerateConfig()
	hasherDal := hasher.Hasher{}
	hasher := hasher.NewHasher(&hasherDal)
	kvDal := kvstore.KVDAL{Store: make(map[string]string)}
	kvStore := kvstore.NewStore(&kvDal, hasher)

	router := router.CreateRouter(kvStore, hasher)

	addr := ":13800"
	srv := &http.Server{
		Handler: router,
		Addr:    addr,
		//Unsure of the timeouts he would want for this
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting on %s\n", addr)
	log.Fatal(srv.ListenAndServe())

}
