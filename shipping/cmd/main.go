package main

import (
	"log"

	"github.com/sofialeals/microservices/shipping/config"
	grpcadapter "github.com/sofialeals/microservices/shipping/internal/adapters/grpc"
	"github.com/sofialeals/microservices/shipping/internal/application/core/api"
)

func main() {
	application := api.NewApplication()
	grpcAdapter := grpcadapter.NewAdapter(application, config.GetApplicationPort())
	if err := grpcAdapter.Run(); err != nil {
		log.Fatalf("failed to start shipping service: %v", err)
	}
}
