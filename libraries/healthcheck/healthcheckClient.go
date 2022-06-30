package healthcheck

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"wartech-studio.com/monster-reacher/libraries/config"
	"wartech-studio.com/monster-reacher/libraries/healthcheck/services/services_discovery"
)

type HealthCheckClient interface {
	Start(name string, host string)
}

type healthCheckClient struct {
	name string
	host string
}

func NewHealthCheckClient() HealthCheckClient {
	return &healthCheckClient{}
}

func (hc *healthCheckClient) Start(name string, host string) {
	hc.name = name
	hc.host = host
	serviceDiscoveryHost := fmt.Sprintf("%s:%d",
		config.WartechConfig().Services["services-discovery"].Hosts[0],
		config.WartechConfig().Services["services-discovery"].Ports[0])
	cc, err := grpc.Dial(serviceDiscoveryHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect: ", err)
		for {
			cc.Connect()
			if state := cc.GetState(); state == connectivity.Ready {
				break
			}
			time.Sleep(3 * time.Second)
		}
	}

	defer cc.Close()

	if state := cc.GetState(); state != connectivity.Ready {
		for {
			cc.Connect()
			if state := cc.GetState(); state == connectivity.Ready {
				break
			}
			time.Sleep(3 * time.Second)
		}
	}

	stream, token, err := hc.registerToServicesDiscovery(cc)
	if err != nil {
		for {
			time.Sleep(3 * time.Second)
			stream, token, err = hc.registerToServicesDiscovery(cc)
			if err == nil {
				break
			}
		}
	}

	for {
		time.Sleep(3 * time.Second)
		if stream == nil {
			cc.Connect()
			stream, token, _ = hc.registerToServicesDiscovery(cc)
			continue
		}
		stream.Send(&services_discovery.HealthCheckRequest{Token: token})
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("health check is error %s", err.Error())
			time.Sleep(3 * time.Second)
			cc.Connect()
			stream, token, _ = hc.registerToServicesDiscovery(cc)
			continue
		}

		if !msg.Success {
			log.Printf("health check is fail %s", msg.Message)
			continue
		}
	}

}

func (hc *healthCheckClient) registerToServicesDiscovery(cc *grpc.ClientConn) (services_discovery.ServicesDiscovery_HealthCheckClient, string, error) {
	if cc.GetState() != connectivity.Ready {
		return nil, "", nil
	}

	c := services_discovery.NewServicesDiscoveryClient(cc)

	registerCtx, registerCtxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer registerCtxCancel()
	res, err := c.Register(registerCtx, &services_discovery.RegisterRequest{
		Service: hc.name,
		Host:    hc.host,
	})

	if err != nil {
		panic(err)
	}

	log.Printf("register to services-discovery is success and get token %s", res.GetToken())

	stream, err := c.HealthCheck(context.Background())

	return stream, res.Token, err
}
