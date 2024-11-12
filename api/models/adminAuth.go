package models

type AdminLoginRequest struct {
	Login    string `json:"phone"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Phone        string `json:"phone"`
	Id           string `json:"id"`
}

type AdminAuthInfo struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type AdminRegisterRequest struct {
	MobilePhone string `json:"mobile_phone"`
}

type AdminRegisterConfRequest struct {
	MobilePhone string `json:"mail"`
	Otp         string `json:"otp"`
	User        *User  `json:"user"`
}

type AdminLoginPhoneConfirmRequest struct {
	MobilePhone string `json:"mobile_phone"`
	SmsCode     string `json:"smscode"`
}
