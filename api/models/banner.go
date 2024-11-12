package models

type Banner struct {
	ImageUrl   string `json:"image_url"`
	Created_at string `json:"created_at,omitempty"`
}

type CreateBanner struct {
	ImageUrl string `json:"image_url"`
}

type DeleteBanner struct {
	Id string `json:"id"`
}

type GetAllBannerRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllBannerResponse struct {
	Banners []Banner `json:"branches"`
	Count   int64    `json:"count"`
}
