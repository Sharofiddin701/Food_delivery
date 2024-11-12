package models

type OrderItem struct {
	Id         string  `json:"id"`
	ProductId  string  `json:"product_id"`
	OrderId    string  `json:"order_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type CreateOrderItem struct {
	ProductId  string  `json:"product_id"`
	OrderId    string  `json:"order_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice int64   `json:"total_price"`
}

type UpdateOrderItem struct {
	ProductId string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type GetOrderItem struct {
	Id        string  `json:"id"`
	OrderId   string  `json:"order_id"`
	ProductId string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type GetAllOrderItemsRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllOrderItemsResponse struct {
	OrderItems []OrderItem `json:"order_items"`
	Count      int64       `json:"count"`
}

type OrderItemsGetListResponse struct {
	Count int          `json:"count"`
	Items []*OrderItem `json:"order_items"`
}

type SwaggerOrderItems struct {
	ProductId string  `json:"product_id,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
}
