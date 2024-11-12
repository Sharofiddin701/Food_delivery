package models

type Branch struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CreateBranch struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UpdateBranch struct {
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type GetBranch struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type GetAllBranchesRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllBranchesResponse struct {
	Branches []Branch `json:"branches"`
	Count    int64    `json:"count"`
}
