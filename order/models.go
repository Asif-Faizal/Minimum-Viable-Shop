package order

import "time"

type Order struct {
	ID         string          `json:"id"`
	CreatedAt  time.Time       `json:"createdAt"`
	AccountID  string          `json:"accountId"`
	TotalPrice float64         `json:"totalPrice"`
	Products   []*OrderProduct `json:"products"`
}

type OrderProduct struct {
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
