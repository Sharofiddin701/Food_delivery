package models

type Payment struct {
	Id            string `json:"id"`
	UserId        string `json:"user_id"`
	OrderId       string `json:"order_id"`
	IsPaid        bool   `json:"is_paid"`
	PaymentMethod string `json:"payment_method"`
	CreatedAt     string `json:"created_at"`
}

type CreatePayment struct {
	UserId        string `json:"user_id"`
	OrderId       string `json:"order_id"`
	IsPaid        bool   `json:"is_paid"`
	PaymentMethod string `json:"payment_method"`
}

type UpdatePayment struct {
	IsPaid        bool   `json:"is_paid"`
	PaymentMethod string `json:"payment_method"`
}

type GetPayment struct {
	Id            string `json:"id"`
	UserId        string `json:"user_id"`
	OrderId       string `json:"order_id"`
	IsPaid        bool   `json:"is_paid"`
	PaymentMethod string `json:"payment_method"`
	CreatedAt     string `json:"created_at"`
}

type GetAllPaymentsRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllPaymentsResponse struct {
	Payments []Payment `json:"payments"`
	Count    int64     `json:"count"`
}
