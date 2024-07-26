package handler

type RegisterResponse struct {
	Username    string `json:"username" form:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required"`
}

type LoginResponse struct {
	Username string `json:"username" form:"username" validate:"required"`
	Token    any    `json:"token"`
}

type UserInfo struct {
	Username    string `json:"username" form:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"required"`
	Email       string `json:"email" form:"email" validate:"required"`
	Role        string `json:"role" form:"role"`
}

type DashboardResponse struct {
	TotalUser         int `json:"total_user"`
	TotalUserBaru     int `json:"total_new_user"`
	TotalUserActive   int `json:"total_active_user"`
	TotalUserInactive int `json:"total_inactive_user"`
}
