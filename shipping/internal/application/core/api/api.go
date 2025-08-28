package api

import (
	"math"

	"github.com/sofialeals/microservices/shipping/internal/application/core/domain"
	"github.com/sofialeals/microservices/shipping/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct{}

func NewApplication() ports.APIPort {
	return &Application{}
}

func (a *Application) Quote(orderID int64, items []domain.Item) (domain.DeliveryQuote, error) {
	var totalUnits int32
	for _, it := range items {
		if it.Quantity <= 0 || it.ProductCode == "" {
			return domain.DeliveryQuote{}, status.Error(codes.InvalidArgument, "invalid item: product_code and positive quantity are required")
		}
		totalUnits += it.Quantity
	}
	if totalUnits < 0 {
		return domain.DeliveryQuote{}, status.Error(codes.InvalidArgument, "invalid total units")
	}
	days := int32(1)
	if totalUnits > 0 {
		days = int32(math.Ceil(float64(totalUnits) / 5.0))
		if days < 1 {
			days = 1
		}
	}
	return domain.DeliveryQuote{OrderID: orderID, EstimatedDays: days}, nil
}
