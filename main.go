package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/colbyleiske/cse138_assignment2/config"
	"github.com/colbyleiske/cse138_assignment2/hasher"
	"github.com/colbyleiske/cse138_assignment2/kvstore"
	"github.com/colbyleiske/cse138_assignment2/router"
)

func main() {
	config.GenerateConfig()

	ringDAL := hasher.NewRing()
	ring := hasher.NewRingStore(ringDAL)

	for _, serverIP := range config.Config.Servers {
		ring.DAL().AddServer(serverIP)
	}

	kvDal := kvstore.KVDAL{Store: make(map[string]string)}
	kvStore := kvstore.NewStore(&kvDal, ring)

	router := router.CreateRouter(kvStore, ring)

	addr := config.Config.Address
	srv := &http.Server{
		Handler: router,
		Addr:    addr,
		//Unsure of the timeouts he would want for this
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting on %s\n", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)

}
