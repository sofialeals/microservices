package ports

import "github.com/sofialeals/microservices/shipping/internal/application/core/domain"

type APIPort interface {
	Quote(orderID int64, items []domain.Item) (domain.DeliveryQuote, error)
}
