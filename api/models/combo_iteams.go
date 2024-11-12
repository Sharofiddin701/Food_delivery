package models

type ComboItem struct {
	Id         string  `json:"id"`
	ProductId  string  `json:"product_id"`
	ComboId    string  `json:"combo_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
}

type CreateComboItem struct {
	ProductId  string  `json:"product_id"`
	ComboId    string  `json:"combo_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice int64   `json:"total_price"`
}

type UpdateComboItem struct {
	ProductId string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type GetComboItem struct {
	Id        string  `json:"id"`
	ComboId   string  `json:"combo_id"`
	ProductId string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type GetAllComboItemsRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllComboItemsResponse struct {
	ComboItems []ComboItem `json:"combo_items"`
	Count      int64       `json:"count"`
}

type ComboItemsGetListResponse struct {
	Count int          `json:"count"`
	Items []*ComboItem `json:"combo_items"`
}
