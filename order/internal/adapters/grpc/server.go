package grpc

import (
	"context"
	"net"
	"strconv"

	"log"

	"github.com/sofialeals/microservices-proto/golang/order"
	"github.com/sofialeals/microservices/order/config"
	"github.com/sofialeals/microservices/order/internal/application/core/domain"
	"github.com/sofialeals/microservices/order/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	api  ports.APIPort
	port int
	order.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem
	for _, orderItem := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	newOrder := domain.NewOrder(int64(request.CostumerId), orderItems)

	result, err := a.api.PlaceOrder(newOrder)
	code := status.Code(err)

	if code == codes.InvalidArgument {
		return nil, err
	} else if err != nil {
		return nil, status.New(
			codes.Internal,
			"failed to place order: "+err.Error(),
		).Err()
	}

	return &order.CreateOrderResponse{OrderId: int32(result.ID)}, nil
}

func (a Adapter) Run() {
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}
	grpcServer := grpc.NewServer()
	order.RegisterOrderServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port")
	}
}
