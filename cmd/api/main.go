package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iruiz/gin-blog-api/internal/applog"
	"github.com/iruiz/gin-blog-api/internal/config"
	"github.com/iruiz/gin-blog-api/internal/database"
	"github.com/iruiz/gin-blog-api/internal/metrics"
	"github.com/iruiz/gin-blog-api/internal/routes"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)
	observability := metrics.NewMetrics(nil)

	db, err := database.New(cfg.DBPath, observability)
	if err != nil {
		applog.Logger.Fatalf("error inicializando BD: %v", err)
	}

	refresherCtx, cancelRefresher := context.WithCancel(context.Background())
	defer cancelRefresher()
	database.StartGaugeRefresher(refresherCtx, db, observability, 30*time.Second)

	router := routes.Setup(db, observability)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	applog.Logger.Printf("Servidor escuchando en :%s (modo %s, versión %s)", cfg.Port, cfg.GinMode, cfg.AppVersion)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			applog.Logger.Fatalf("error arrancando servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	applog.Logger.Println("Apagando servidor...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		applog.Logger.Fatalf("error apagando servidor: %v", err)
	}
}
