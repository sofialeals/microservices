package api

import (
	"github.com/sofialeals/microservices/order/internal/application/core/domain"
	"github.com/sofialeals/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	totalItens := 0
	for _, item := range order.OrderItems {
		totalItens += int(item.Quantity)
	}

	if totalItens > 50 {
		order.Status = "Canceled"
		_ = a.db.Save(&order)
		return domain.Order{}, status.Errorf(
			codes.InvalidArgument,
			"Order with %d items exceeds the maximum allowed of 50", totalItens,
		)
	}

	order.Status = "Paid"
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		order.Status = "Canceled"
		_ = a.db.Save(&order)
		return domain.Order{}, paymentErr
	}

	return order, nil
}
