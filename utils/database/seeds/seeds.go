package seeds

import (
	"e-ticketing-gin/utils/database/seed"
	"gorm.io/gorm"
)

func All() []seed.Seed {
	var seeds []seed.Seed = []seed.Seed{
		{
			Name: "Create Admin",
			Run: func(db *gorm.DB) error {
				return CreateUser(db, "Admin", "admin@gmail.com", "081317394655")
			},
		},
	}

	return seeds
}
