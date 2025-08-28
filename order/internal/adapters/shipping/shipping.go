package shipping_test

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/sofialeals/microservices-proto/golang/shipping"
	"github.com/sofialeals/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client shipping.ShippingClient
}

func NewAdapter(shippingServiceURL string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
				grpc_retry.WithMax(5),
				grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
			),
		),
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(shippingServiceURL, opts...)
	if err != nil {
		return nil, err
	}
	return &Adapter{client: shipping.NewShippingClient(conn)}, nil
}

func (a *Adapter) Schedule(order *domain.Order) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var items []*shipping.Item
	for _, it := range order.OrderItems {
		items = append(items, &shipping.Item{ProductCode: it.ProductCode, Quantity: it.Quantity})
	}
	res, err := a.client.Create(ctx, &shipping.CreateShippingRequest{
		OrderId: order.ID,
		Items:   items,
	})
	if err != nil {
		return 0, err
	}
	return res.EstimatedDays, nil
}
