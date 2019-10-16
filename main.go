package main

import (
	"log"
	"net/http"

	"./kvstore"
	"./router"
)

func main() {
	dal := kvstore.KVDAL{Store: make(map[string]string)}
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
