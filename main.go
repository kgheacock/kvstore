package main

import (
	"log"
	"net/http"

	"github.com/colbyleiske/cse138_assignment2/kvstore"
	"github.com/colbyleiske/cse138_assignment2/router"
)

func main() {
	dal := kvstore.KVDAL{}
	kvStore := kvstore.NewStore(&dal)

	router := router.CreateRouter(kvStore)

	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:13800",
		//Unsure of the timeouts he would want for this
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting up...")
	log.Fatal(srv.ListenAndServe())

}
