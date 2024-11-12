package models

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateCategory struct {
	Name string `json:"name"`
}

type UpdateCategory struct {
	Name string `json:"name"`
}

type GetCategory struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetAllCategoriesRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllCategoriesResponse struct {
	Categories []Category `json:"categories"`
	Count      int64      `json:"count"`
}
