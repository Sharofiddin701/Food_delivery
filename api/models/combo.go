package models

type GetAllCombosRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllCombosResponse struct {
	Combos []Combo `json:"combos"`
	Count  int64   `json:"count"`
}

type Combo struct {
	Id          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Price       float64     `json:"price,omitempty"`
	TotalPrice  float64     `json:"total_price,omitempty"`
	Status      string      `json:"status,omitempty"`
	CreatedAt   string      `json:"created_at,omitempty"`
	UpdatedAt   string      `json:"updated_at,omitempty"`
	ComboItems  []ComboItem `json:"combo_items,omitempty"`
}

type ComboCreate struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
}

type SwaggerComboCreate struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
}

type ComboUpdate struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Price       float64     `json:"price"`
	ComboItems  []ComboItem `json:"combo_items,omitempty"`
}

type ComboUpdateS struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Price       float64           `json:"price"`
	ComboItems  []UpdateComboItem `json:"combo_items,omitempty"`
}

type ComboPrimaryKey struct {
	Id string `json:"id"`
}

type ComboGetListRequest struct {
	Search string `json:"search,omitempty"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type ComboGetListResponse struct {
	Count int      `json:"count"`
	Combo []*Combo `json:"combo"`
}

type ComboCreateRequest struct {
	Combo Combo       `json:"combo"`
	Items []ComboItem `json:"items"`
}

type SwaggerComboCreateRequest struct {
	Combo SwaggerComboCreate  `json:"combo"`
	Items []SwaggerComboItems `json:"items"`
}

type SwaggerComboItems struct {
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
