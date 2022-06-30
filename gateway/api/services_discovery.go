package api

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"wartech-studio.com/monster-reacher/gateway/services/services_discovery"
)

type ServicesDiscovery interface {
	Start(host string)
	GetServiceInfo(name string) *services_discovery.ServiceInfo
}

type servicesDiscovery struct {
	services []*services_discovery.ServiceInfo
}

func NewServicesDiscovery() ServicesDiscovery {
	return &servicesDiscovery{}
}

func (sd *servicesDiscovery) Start(host string) {
	for {

		cc, err := grpc.Dial("localhost:20531", grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Println("fetch api error " + err.Error())
			continue
		}

		c := services_discovery.NewServicesDiscoveryClient(cc)

		res, err := c.GetServices(context.Background(), &services_discovery.GetServicesRequest{})

		if err != nil {
			log.Println("get services error " + err.Error())
			continue
		}

		for _, service := range res.GetServices() {
			res, err := c.CheckServiceIsOnline(context.Background(), &services_discovery.CheckServiceIsOnlineRequest{
				Name: service.GetName(),
			})
			if err != nil {
				log.Println("check service is online error " + err.Error())
				continue
			}

			if res.GetInfo().GetIsOnline() {
				sd.services = append(sd.services, res.GetInfo())
			}
		}

		log.Println(sd.services)

		time.Sleep(5 * time.Second)
	}
}

func (sd *servicesDiscovery) GetServiceInfo(name string) *services_discovery.ServiceInfo {
	for _, s := range sd.services {
		if s.GetName() == name {
			return s
		}
	}

	return nil
}
