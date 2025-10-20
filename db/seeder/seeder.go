package seeder

import (
	"log"
	"pura-agung-kertajaya-backend/internal/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		log.Println("Seeder: users already exist, skipping...")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("rahasia"), bcrypt.DefaultCost)

	users := []entity.User{
		{
			Name:     "Administrator",
			Email:    "admin@puraagungkertajaya.com",
			Password: string(password),
		},
		{
			Name:     "Joko",
			Email:    "joko@example.com",
			Password: string(password),
		},
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Seeder error: %v", err)
	}
	log.Println("Seeder: users table seeded successfully!")
}
