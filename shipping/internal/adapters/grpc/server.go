package grpc

import (
	"context"
	"net"
	"strconv"

	"github.com/sofialeals/microservices-proto/golang/shipping"
	"github.com/sofialeals/microservices/shipping/config"
	"github.com/sofialeals/microservices/shipping/internal/application/core/domain"
	"github.com/sofialeals/microservices/shipping/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a *Adapter) Create(ctx context.Context, req *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	var items []domain.Item
	for _, it := range req.Items {
		items = append(items, domain.Item{ProductCode: it.ProductCode, Quantity: it.Quantity})
	}
	quote, err := a.api.Quote(req.OrderId, items)
	if err != nil {
		return nil, err
	}
	return &shipping.CreateShippingResponse{
		OrderId:       quote.OrderID,
		EstimatedDays: quote.EstimatedDays,
	}, nil
}

func (a *Adapter) Run() error {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	shipping.RegisterShippingServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}
	return grpcServer.Serve(listen)
}
