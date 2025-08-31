package main

// Order представляет начальную модель заказа
type Order struct {
	OrderUID   string `json:"order_uid"`
	CustomerID string `json:"customer_id"`
}