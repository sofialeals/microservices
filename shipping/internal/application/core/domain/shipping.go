package domain

type Item struct {
	ProductCode string
	Quantity    int32
}

type DeliveryQuote struct {
	OrderID       int64
	EstimatedDays int32
}
