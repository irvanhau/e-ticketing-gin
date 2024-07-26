package handler

type RegisterInput struct {
	Username    string `json:"username" form:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required"`
	Password    string `json:"password" form:"password" validate:"required"`
}

type LoginInput struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type ForgetPasswordInput struct {
	Username string `json:"username" form:"username" validate:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" form:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm" validate:"required"`
}

type UpdateProfile struct {
	Username    string `json:"username" form:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required"`
}

type RefreshTokenInput struct {
	Token string `json:"access_token" form:"access_token"`
}
