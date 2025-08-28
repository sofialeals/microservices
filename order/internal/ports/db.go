package ports

import "github.com/sofialeals/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
	ProductCodesExist(codes []string) ([]string, error)
}
