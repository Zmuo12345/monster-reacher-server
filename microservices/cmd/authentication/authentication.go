package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"wartech-studio.com/monster-reacher/microservices/services/authentication"

	"wartech-studio.com/monster-reacher/libraries/config"
	"wartech-studio.com/monster-reacher/libraries/healthcheck"
)

const SERVICES_NAME = "authentication"

var listenHost = fmt.Sprintf("%s:%d",
	config.WartechConfig().Services[SERVICES_NAME].Hosts[0],
	config.WartechConfig().Services[SERVICES_NAME].Ports[0])

func main() {
	server := grpc.NewServer()
	serviceDiscoveryHost := fmt.Sprintf("%s:%d",
		config.WartechConfig().Services["services-discovery"].Hosts[0],
		config.WartechConfig().Services["services-discovery"].Ports[0])
	healthchecker := healthcheck.NewHealthCheckClient()
	go healthchecker.Start(SERVICES_NAME, serviceDiscoveryHost)
	listener, err := net.Listen("tcp", listenHost)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	authentication.RegisterAuthenticationServer(server, authentication.NewAuthenticationServer())
	log.Println("gRPC server listening on " + listenHost)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
