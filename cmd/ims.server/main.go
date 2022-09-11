package main

import (
	"cmd/ims.server/pkg/adding"
	"cmd/ims.server/pkg/config"
	"cmd/ims.server/pkg/http/rest"
	"cmd/ims.server/pkg/listing"
	"cmd/ims.server/pkg/storage/postgres"
	"cmd/ims.server/pkg/updating"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to read config file:", err)
	}

	var adder adding.Service
	var lister listing.Service
	var updater updating.Service
	// error handling omitted for simplicity
	s, _ := postgres.NewRepository()

	adder = adding.NewService(s)
	lister = listing.NewService(s)
	updater = updating.NewService(s)

	router := rest.Handler(adder, lister, updater)
	fmt.Println("The ims server is on tap now: http://localhost:8080")
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
