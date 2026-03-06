package initializers

import (
	"Gin/models"
)

// AutoMigrate runs GORM automigrations for all models.
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Company{},
		&models.House{},
	)
}
