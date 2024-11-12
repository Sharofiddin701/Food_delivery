package models

type GetAllOrdersRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllOrdersResponse struct {
	Orders []Order `json:"orders"`
	Count  int64   `json:"count"`
}

type ChangeStatus struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type GetOrderStatus struct {
	Status string `json:"status"`
}

type Order struct {
	Id             string      `json:"id,omitempty"`
	UserId         string      `json:"user_id,omitempty"`
	TotalPrice     float64     `json:"total_price,omitempty"`
	Longitude      float64     `json:"longitude"`
	Latitude       float64     `json:"latitude"`
	AddressName    string      `json:"address_name"`
	Status         string      `json:"status,omitempty"`
	DeliveryStatus string      `json:"delivery_status"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	OrderItems     []OrderItem `json:"order_items,omitempty"`
}

type OrderCreate struct {
	UserId      string  `json:"user_id,omitempty"`
	Status      string  `json:"status"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
	AddressName string  `json:"address_name"`
}

type SwaggerOrderCreate struct {
	UserId         string  `json:"user_id,omitempty"`
	DeliveryStatus string  `json:"delivery_status"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	AddressName    string  `json:"address_name"`
}

type OrderUpdate struct {
	Longitude   float64     `json:"longitude"`
	Latitude    float64     `json:"latitude"`
	AddressName string      `json:"address_name"`
	TotalPrice  float64     `json:"total_price"`
	Status      string      `json:"status"`
	OrderItems  []OrderItem `json:"order_items,omitempty"`
}

type OrderUpdateS struct {
	Longitude   float64           `json:"longitude"`
	Latitude    float64           `json:"latitude"`
	AddressName string            `json:"address_name"`
	TotalPrice  float64           `json:"total_price"`
	Status      string            `json:"status"`
	OrderItems  []UpdateOrderItem `json:"order_items,omitempty"`
}

type OrderPrimaryKey struct {
	Id string `json:"id"`
}

type OrderGetListRequest struct {
	UserId string `json:"user_id,omitempty"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type OrderGetListResponse struct {
	Count int      `json:"count"`
	Order []*Order `json:"order"`
}

type OrderCreateRequest struct {
	Order Order       `json:"order"`
	Items []OrderItem `json:"items"`
}

type SwaggerOrderCreateRequest struct {
	Order SwaggerOrderCreate  `json:"order"`
	Items []SwaggerOrderItems `json:"items"`
}

type PatchOrderStatusRequest struct {
	Status string `json:"status"`
}
