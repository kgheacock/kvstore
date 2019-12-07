package main

import (
	"context"
	"log"
	"math/rand"
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
	rand.Seed(4091999) //gonna use my birthdate for deterministic testing

	config.GenerateConfig()

	ringDAL := hasher.NewRing()
	ring := hasher.NewRingStore(ringDAL)

	for _, shard := range config.Config.Shards {
		ringDAL.AddShard(shard)
	}

	kvDal := kvstore.KVDAL{Store: make(map[string]kvstore.StoredValue)}
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
