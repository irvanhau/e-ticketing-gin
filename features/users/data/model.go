package data

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	*gorm.Model
	Username    string `gorm:"column:username;type:varchar(255);not null"`
	Email       string `gorm:"column:email;type:varchar(255);not null"`
	PhoneNumber string `gorm:"column:phone_number;type:varchar(255);not null"`
	Password    string `gorm:"column:password;type:varchar(255);not null"`
	IsAdmin     bool   `gorm:"column:is_admin;type:bool;not null"`
	Status      bool   `gorm:"column:status;type:bool;not null"`
}

type UserResetPass struct {
	*gorm.Model
	Username  string    `gorm:"column:username;type:varchar(255);not null"`
	Code      string    `gorm:"column:code;type:varchar(255);not null"`
	ExpiredAt time.Time `gorm:"column:expired_at;type:timestamp;not null"`
}

type UserVerification struct {
	Username  string    `gorm:"column:username;type:varchar(255);not null"`
	Code      string    `gorm:"column:code;type:varchar(255);not null"`
	ExpiredAt time.Time `gorm:"column:expired_at;type:timestamp;not null"`
}
