package models

type Product struct {
	Id          string  `json:"id"`
	CategoryId  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateProduct struct {
	CategoryId  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

type UpdateProduct struct {
	CategoryId  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

type GetProduct struct {
	Id          string  `json:"id"`
	CategoryId  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type GetAllProductsRequest struct {
	CategoryId string `json:"category_id"`
	Search     string `json:"search"`
	Page       uint64 `json:"page"`
	Limit      uint64 `json:"limit"`
}

type GetAllProductsResponse struct {
	Products []Product `json:"products"`
	Count    int64     `json:"count"`
}
