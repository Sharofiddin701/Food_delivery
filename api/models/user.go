package models

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Sex        string `json:"sex"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Created_at string `json:"created_at,omitempty"`
	Updated_at string `json:"updated_at,omitempty"`
}

type CreateUser struct {
	Name  string `json:"name"`
	Sex        string `json:"sex"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UpdateUser struct {
	Name  string `json:"name"`
	Sex        string `json:"sex"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type GetUser struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Sex        string `json:"sex"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type GetAllUsersRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllUsersResponse struct {
	Users []User `json:"users"`
	Count int64  `json:"count"`
}
