package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/config"
	"github.com/iruiz/gin-blog-api/internal/database"
	"github.com/iruiz/gin-blog-api/internal/routes"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("error inicializando BD: %v", err)
	}

	router := routes.Setup(db)

	log.Printf("Servidor escuchando en :%s (modo %s, versión %s)", cfg.Port, cfg.GinMode, cfg.AppVersion)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("error arrancando servidor: %v", err)
	}
}
