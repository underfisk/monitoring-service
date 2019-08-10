package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/underfisk/monitoring-service/core"
)

/**
	This service collects metrics related with multiple servers/instances/apps
	Also will be used to provide logs monitoring, service monitoring and frontend application
	interface.
	In order to manage the dashboard (visually) we'll be using Vue.js as a spa
	and authorization based using JWT

	Support:
		- gRPC (For fast communication and create nestjs client for this)
		- REST (http support)
 */
func main () {
	print("Running monitoring Service on port: 4000 for now")



	http.HandleFunc("/bar",

	log.Fatal(http.ListenAndServe(":4000", nil))
}

