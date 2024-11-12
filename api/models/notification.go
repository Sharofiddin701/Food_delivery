package models

type Notification struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type CreateNotification struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}

type UpdateNotification struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}

type GetNotification struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type GetAllNotificationsRequest struct {
	Search string `json:"search"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type GetAllNotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
	Count         int64          `json:"count"`
}
