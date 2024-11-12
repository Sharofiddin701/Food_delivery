package models

type DeliveryHistory struct {
	Id          string  `json:"id"`
	CourierId   string  `json:"courier_id"`
	OrderId     string  `json:"order_id"`
	Earnings    float64 `json:"earnings"`
	DeliveredAt string  `json:"delivered_at"`
}

type CreateDeliveryHistory struct {
	CourierId   string  `json:"courier_id"`
	OrderId     string  `json:"order_id"`
	Earnings    float64 `json:"earnings"`
	DeliveredAt string  `json:"delivered_at"`
}

type UpdateDeliveryHistory struct {
	CourierId   string  `json:"courier_id"`
	OrderId     string  `json:"order_id"`
	Earnings    float64 `json:"earnings"`
	DeliveredAt string  `json:"delivered_at"`
}

type GetDeliveryHistory struct {
	Id          string  `json:"id"`
	CourierId   string  `json:"courier_id"`
	OrderId     string  `json:"order_id"`
	Earnings    float64 `json:"earnings"`
	DeliveredAt string  `json:"delivered_at"`
}

type GetAllDeliveryHistoriesRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllDeliveryHistoriesResponse struct {
	DeliveryHistories []DeliveryHistory `json:"delivery_histories"`
	Count             int64             `json:"count"`
}
