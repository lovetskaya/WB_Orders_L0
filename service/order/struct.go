package order

import (
	"time"
)

type Delivery struct {
	Name    string `json:"delivery_name"`
	Phone   string `json:"delivery_phone"`
	Zip     string `json:"delivery_zip"`
	City    string `json:"delivery_city"`
	Address string `json:"delivery_address"`
	Region  string `json:"delivery_region"`
	Email   string `json:"delivery_email"`
}

type Payment struct {
	Transaction    string    `json:"payment_transaction"`
	RequestID      string    `json:"payment_request_id"`
	Currency       string    `json:"payment_currency"`
	Provider       string    `json:"payment_provider"`
	Amount         float64   `json:"payment_amount"`
	PaymentDt      time.Time `json:"payment_payment_dt"`
	Bank           string    `json:"payment_bank"`
	DeliveryCost   float64   `json:"payment_delivery_cost"`
	GoodsTotal     float64   `json:"payment_goods_total"`
	CustomFee      float64   `json:"payment_custom_fee"`
}

type Item struct {
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float64 `json:"price"`
	Rid         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SMID              int      `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}
