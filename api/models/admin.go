package models

type Admin struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
}

type CreateAdmin struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UpdateAdmin struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type GetAdmin struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type GetAllAdminsRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllAdminsResponse struct {
	Admins []Admin `json:"Admins"`
	Count  int64   `json:"count"`
}
