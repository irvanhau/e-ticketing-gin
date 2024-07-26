package database

import (
	"e-ticketing-gin/configs"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(c *configs.ProgramConfig) *gorm.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.DBUser, c.DBPass, c.DBHost, c.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logrus.Error("Terjadi kesalahan pada database error: ", err.Error())
		return nil
	}

	return db
}
