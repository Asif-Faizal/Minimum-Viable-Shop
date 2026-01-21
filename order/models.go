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
	OrderID            string  `json:"orderId"`
	ProductID          string  `json:"productId"`
	ProductName        string  `json:"productName"`
	ProductDescription string  `json:"productDescription"`
	Price              float64 `json:"price"`
	Quantity           int     `json:"quantity"`
}
