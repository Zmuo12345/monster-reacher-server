package services_discovery

import (
	context "context"
	"log"

	"github.com/google/uuid"
)

type servicesDiscoveryServer struct {
	services map[string]*ServiceInfo
}

func NewServicesDiscoveryServer() ServicesDiscoveryServer {
	return &servicesDiscoveryServer{
		services: make(map[string]*ServiceInfo),
	}
}

func (sd *servicesDiscoveryServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	for k := range sd.services {
		if sd.services[k].Host == req.GetHost() && sd.services[k].Port == req.GetPort() {
			if sd.services[k].Name != req.GetService() {
				log.Printf("the service on host %s port %s is switching from %s to %s\n", req.GetHost(), req.GetPort(), sd.services[k].Name, req.GetService())
			}
			sd.services[k].Name = req.GetService()
			log.Printf("the service %s on host %s port %s was reconnect \n", req.GetService(), req.GetHost(), req.GetPort())
			return &RegisterResponse{Token: k}, nil
		}
	}
	uuid := uuid.New()
	sd.services[uuid.String()] = &ServiceInfo{
		Name: req.GetService(),
		Host: req.GetHost(),
		Port: req.GetPort(),
	}
	log.Printf("the service %s on host %s port %s was register is success \n", req.GetService(), req.GetHost(), req.GetPort())
	return &RegisterResponse{Token: string(uuid.String())}, nil
}

func (sd *servicesDiscoveryServer) HealthCheck(stream ServicesDiscovery_HealthCheckServer) error {
	token := ""
	for {
		msg, err := stream.Recv()
		if _, found := sd.services[msg.GetToken()]; !found {
			service := sd.services[token]
			service.IsOnline = false
			log.Printf("the service %s on host %s port %s is closed %s\n", service.Name, service.Host, service.Port, err.Error())
			break
		}
		service := sd.services[msg.GetToken()]
		if err != nil {
			log.Printf("the service %s on host %s port %s is closed %s\n", service.Name, service.Host, service.Port, err.Error())
			break
		}
		stream.Send(&HealthCheckResponse{Success: true, Message: "ok"})
		token = msg.GetToken()
		service.IsOnline = true
		log.Printf("the service %s on host %s port %s health check is ok \n", service.Name, service.Host, service.Port)
	}
	return nil
}
func (sd *servicesDiscoveryServer) CheckServiceIsOnline(ctx context.Context, req *CheckServiceIsOnlineRequest) (*CheckServiceIsOnlineResponse, error) {
	for k := range sd.services {
		if sd.services[k].GetName() == req.GetName() && sd.services[k].GetIsOnline() {
			return &CheckServiceIsOnlineResponse{
				IsOnline: sd.services[k].GetIsOnline(),
			}, nil
		}
	}
	return &CheckServiceIsOnlineResponse{IsOnline: false}, nil
}
func (sd *servicesDiscoveryServer) GetServices(context.Context, *GetServicesRequest) (*GetServicesresponse, error) {
	services := []*ServiceInfo{}

	for k := range sd.services {
		services = append(services, sd.services[k])
	}
	return &GetServicesresponse{
		Services: services,
	}, nil
}
func (*servicesDiscoveryServer) GatewaySocket(ServicesDiscovery_GatewaySocketServer) error {
	return nil
}
func (*servicesDiscoveryServer) mustEmbedUnimplementedServicesDiscoveryServer() {}