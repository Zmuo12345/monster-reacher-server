protoc *.proto --go_out=../microservices/services --go-grpc_out=../microservices/services
protoc *.proto --go_out=../gateway/services --go-grpc_out=../gateway/services
protoc services_discovery.proto --go_out=../libraries/healthcheck --go-grpc_out=../libraries/healthcheck
protoc services_discovery.proto --go_out=../services-discovery --go-grpc_out=../services-discovery