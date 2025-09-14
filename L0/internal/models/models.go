package models

import (
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid" fake:"{uuid}"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery" fake:"skip"`
	Payment           Payment   `json:"payment" fake:"skip"`
	Items             []Item    `json:"items" fake:"skip"`
	Locale            string    `json:"locale" fake:"{languageabbreviation}"`
	InternalSignature string    `json:"internal_signature" fake:"skip"`
	CustomerID        string    `json:"customer_id" fake:"{uuid}""`
	DeliveryService   string    `json:"delivery_service" fake:"{company}"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id" fake:"{number:1,100}"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	OrderUID string `json:"-"`
	Name     string `json:"name" fake:"{name}"`
	Phone    string `json:"phone" fake:"{phone}"`
	Zip      string `json:"zip"`
	City     string `json:"city" fake:"{city}"`
	Address  string `json:"address" fake:"{address}"`
	Region   string `json:"region" fake:"{state}"`
	Email    string `json:"email" fake:"{email}"`
}

type Payment struct {
	OrderUID     string  `json:"-" db:"order_uid"`
	Transaction  string  `json:"transaction" fake:"{uuid}"`
	RequestID    string  `json:"request_id" fake:"{uuid}"`
	Currency     string  `json:"currency" fake:"{CurrencyShort}"`
	Provider     string  `json:"provider" fake:"{company}"`
	Amount       float64 `json:"amount" fake:"{price:1,10000}"`
	PaymentDt    int     `json:"payment_dt" fake:"{number:100,1000}"`
	Bank         string  `json:"bank" fake:"{bankname}"`
	DeliveryCost float64 `json:"delivery_cost" fake:"{price:1,10000}"`
	GoodsTotal   float64 `json:"goods_total" fake:"{price:1,10000}"`
	CustomFee    float64 `json:"custom_fee" fake:"{price:1,10000}"`
}

type Item struct {
	ID          int    `json:"-"`
	OrderUID    string `json:"-"`
	ChrtID      int64  `json:"chrt_id"fake:"{number:1,10000}"`
	TrackNumber string `json:"track_number" `
	Price       int    `json:"price" fake:"{number:1000,10000}"`
	Rid         string `json:"rid" fake:"{uuid}"`
	Name        string `json:"name" fake:"{productname}"`
	Sale        int    `json:"sale" fake:"{number:0,100}"`
	Size        string `json:"size" fake:"{number:0,100}"`
	TotalPrice  int    `json:"total_price" fake:"{number:1,10000} `
	NmID        int64  `json:"nm_id" fake:"{number:10000,99999}"`
	Brand       string `json:"brand" fake:"{company}"`
	Status      int    `json:"status" fake:"{number:200,202}"`
}
