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
			Role:     "super",
		},
		{
			Name:     "Jonathan",
			Email:    "jonathan@puraagungkertajaya.com",
			Password: string(password),
			Role:     "super",
		},
		{
			Name:     "Admin Pura",
			Email:    "pura@puraagungkertajaya.com",
			Password: string(password),
			Role:     "pura",
		},
		{
			Name:     "Admin Yayasan",
			Email:    "yayasan@puraagungkertajaya.com",
			Password: string(password),
			Role:     "yayasan",
		},
		{
			Name:     "Admin Pasraman",
			Email:    "pasraman@puraagungkertajaya.com",
			Password: string(password),
			Role:     "pasraman",
		},
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Seeder error: %v", err)
	}
	log.Println("Seeder: users table seeded successfully!")
}
