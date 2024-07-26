package seeds

import (
	"e-ticketing-gin/features/users"
	"e-ticketing-gin/helper/enkrip"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, username, email, phoneNumber string) error {
	var countData int64
	db.Table("users").Where("email = ?", email).Where("username = ?", username).Count(&countData)

	if countData < 1 {
		hash := enkrip.New()
		hashPass, _ := hash.HashPassword("password")
		return db.Create(&users.User{
			Username:    username,
			Email:       email,
			Password:    hashPass,
			PhoneNumber: phoneNumber,
			IsAdmin:     true,
			Status:      true,
		}).Error
	}
	return nil
}
