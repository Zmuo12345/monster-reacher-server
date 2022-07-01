package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"wartech-studio.com/monster-reacher/gateway/api"
	"wartech-studio.com/monster-reacher/libraries/config"
)

const SERVICES_NAME = "gateway"

var listenHost = fmt.Sprintf("%s:%d",
	config.WartechConfig().Services[SERVICES_NAME].Hosts[0],
	config.WartechConfig().Services[SERVICES_NAME].Ports[0])

func main() {
	initServicesDiscovery()
	router := mux.NewRouter().StrictSlash(true)
	api.RegisterHomeApiHandle(router)
	api.RegisterAuthApiHandle(router)
	log.Fatal(http.ListenAndServe(listenHost, router))
}

func initServicesDiscovery() {
	servicesDiscoveryHost := fmt.Sprintf("%s:%d",
		config.WartechConfig().Services["services-discovery"].Hosts[0],
		config.WartechConfig().Services["services-discovery"].Ports[0])
	go api.ServicesDiscoveryCache.Start(servicesDiscoveryHost)
}
