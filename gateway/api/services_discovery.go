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
	CheckRequireServices(serviceNames []string) (map[string]*services_discovery.ServiceInfo, bool)
}

type servicesDiscovery struct {
	services []*services_discovery.ServiceInfo
}

func NewServicesDiscovery() ServicesDiscovery {
	return &servicesDiscovery{}
}

var ServicesDiscoveryCache ServicesDiscovery = NewServicesDiscovery()

func (sd *servicesDiscovery) Start(host string) {
	for {
		time.Sleep(3 * time.Second)
		cc, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

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
			sd.updateServicesInfo(service)
		}

		log.Println(sd.services)
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

func (sd *servicesDiscovery) updateServicesInfo(service *services_discovery.ServiceInfo) {
	for _, s := range sd.services {
		if s.GetHost() == service.GetHost() {
			s.Host = service.GetHost()
			s.Name = service.GetName()
			s.IsOnline = service.GetIsOnline()
			return
		}
	}

	sd.services = append(sd.services, service)
}

func (sd *servicesDiscovery) CheckRequireServices(serviceNames []string) (map[string]*services_discovery.ServiceInfo, bool) {
	requireServices := make(map[string]*services_discovery.ServiceInfo)
	found := 0
	for _, name := range serviceNames {
		for _, service := range sd.services {
			if service.GetName() == name {
				requireServices[name] = service
				found++
				break
			}
		}
	}

	return requireServices, len(serviceNames) == found
}
