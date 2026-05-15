package database

import (
	"fmt"

	"github.com/iruiz/gin-blog-api/internal/metrics"
	"github.com/iruiz/gin-blog-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func New(dbPath string, observability *metrics.Metrics) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
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

	if err := EnableMetrics(db, observability); err != nil {
		return nil, fmt.Errorf("error activando métricas de BD: %w", err)
	}

	return db, nil
}

func EnableMetrics(db *gorm.DB, observability *metrics.Metrics) error {
	if observability == nil {
		return nil
	}

	return newMetricsCallbacks(observability).register(db)
}
