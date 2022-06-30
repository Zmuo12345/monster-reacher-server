cd ../microservices
go run ./cmd/authentication/authentication.go
cd ../microservices
go run ./cmd/profile/profile.go
cd ../gateway
go run .
cd ../services-discovery
go run .
