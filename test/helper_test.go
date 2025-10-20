package test

import (
	"pura-agung-kertajaya-backend/internal/entity"
)

func ClearAll() {
	ClearUsers()
}

func ClearUsers() {
	err := db.Where("id IS NOT NULL").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear users: %+v", err)
	}
}
