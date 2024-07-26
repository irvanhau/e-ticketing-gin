package database

import (
	"e-ticketing-gin/features/users/data"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(data.User{})
	db.AutoMigrate(data.UserResetPass{})
	db.AutoMigrate(data.UserVerification{})
}
