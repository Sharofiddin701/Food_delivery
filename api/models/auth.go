package models

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Phone        string `json:"phone"`
	Id           string `json:"id"`
}

type AuthInfo struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type UserRegisterRequest struct {
	MobilePhone string `json:"mobile_phone"`
}

type UserRegisterConfRequest struct {
	MobilePhone string `json:"mail"`
	Otp         string `json:"otp"`
	User        *User  `json:"user"`
}

type UserLoginPhoneConfirmRequest struct {
	MobilePhone string `json:"mobile_phone"`
	SmsCode     string `json:"smscode"`
}
