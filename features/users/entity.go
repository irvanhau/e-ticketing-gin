package users

import (
	"github.com/gin-gonic/gin"
	"time"
)

type User struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	IsAdmin     bool   `json:"is_admin"`
	Status      bool   `json:"status"`
}

type UserCredential struct {
	Username string         `json:"username"`
	Access   map[string]any `json:"token"`
}

type UserResetPass struct {
	Username  string    `json:"username"`
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UserVerification struct {
	Username  string    `json:"username"`
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UpdateProfile struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type UserDashboard struct {
	TotalUser         int `json:"total_user"`
	TotalNewUser      int `json:"total_new_user"`
	TotalUserActive   int `json:"total_active_user"`
	TotalUserInactive int `json:"total_inactive_user"`
}
type UserHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ForgetPasswordWeb(c *gin.Context)
	ResetPassword(c *gin.Context)
	UpdateProfile(c *gin.Context)
	RefreshToken(c *gin.Context)
	Profile(c *gin.Context)

	GetUsers(c *gin.Context)
	ActivateUser(c *gin.Context)
	DeactivateUser(c *gin.Context)

	UserDashboard(c *gin.Context)
	UserVerification(c *gin.Context)
}

type UserServiceInterface interface {
	Register(newData User) (*User, error)
	Login(username string, password string) (*UserCredential, error)
	ForgetPasswordWeb(username string) error
	TokenResetVerify(code string) (*UserResetPass, error)
	ResetPassword(code, username, password string) error
	UpdateProfile(id int, newData UpdateProfile) (bool, error)
	Profile(id int) (*User, error)

	GetAll() ([]User, error)
	Activate(id int) (bool, error)
	Deactivate(id int) (bool, error)

	UserDashboard() (UserDashboard, error)
	UserVerificationCode(username, email string) error
	UserVerification(code, username string) error
	TokenVerificationResetVerify(code string) (*UserVerification, error)
}

type UserDataInterface interface {
	Register(newData User) (*User, error)
	Login(username, password string) (*User, error)
	GetByID(id int) (User, error)
	GetByUsername(username string) (*User, error)
	InsertCodeReset(username, code string) error
	DeleteCodeReset(code string) error
	GetByCodeReset(code string) (*UserResetPass, error)
	ResetPassword(code, username, password string) error
	UpdateProfile(id int, newData UpdateProfile) (bool, error)
	CheckUsername(username string) bool

	GetAll() ([]User, error)
	Activate(id int) (bool, error)
	Deactivate(id int) (bool, error)

	UserDashboard() (UserDashboard, error)
	InsertCodeVerification(username, code string) error
	DeleteCodeVerification(code string) error
	GetByCodeVerification(code string) (*UserVerification, error)
	UserVerification(code, username string) error
}
