package database

import (
	"fmt"

	"github.com/iruiz/gin-blog-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("error abriendo BD: %w", err)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("error activando foreign_keys: %w", err)
	}

	if err := db.AutoMigrate(&models.Post{}, &models.Comment{}); err != nil {
		return nil, fmt.Errorf("error migrando esquema: %w", err)
	}

	return db, nil
}
