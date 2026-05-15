package database

import (
	"context"
	"time"

	"github.com/iruiz/gin-blog-api/internal/applog"
	"github.com/iruiz/gin-blog-api/internal/metrics"
	"github.com/iruiz/gin-blog-api/internal/models"
	"gorm.io/gorm"
)

type GaugeRefresher struct {
	db            *gorm.DB
	observability *metrics.Metrics
	interval      time.Duration
}

func NewGaugeRefresher(db *gorm.DB, observability *metrics.Metrics, interval time.Duration) *GaugeRefresher {
	return &GaugeRefresher{
		db:            db,
		observability: observability,
		interval:      interval,
	}
}

func StartGaugeRefresher(ctx context.Context, db *gorm.DB, observability *metrics.Metrics, interval time.Duration) {
	if observability == nil {
		return
	}

	NewGaugeRefresher(db, observability, interval).Start(ctx)
}

func (g *GaugeRefresher) Start(ctx context.Context) {
	if g == nil || g.db == nil || g.observability == nil {
		return
	}

	g.refresh(ctx)

	ticker := time.NewTicker(g.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				g.refresh(ctx)
			}
		}
	}()
}

func (g *GaugeRefresher) refresh(ctx context.Context) {
	if err := g.refreshPosts(ctx, true); err != nil {
		applog.Logger.Printf("refresh gauges: failed counting published posts: %v", err)
	}
	if err := g.refreshPosts(ctx, false); err != nil {
		applog.Logger.Printf("refresh gauges: failed counting draft posts: %v", err)
	}
	if err := g.refreshComments(ctx); err != nil {
		applog.Logger.Printf("refresh gauges: failed counting comments: %v", err)
	}
}

func (g *GaugeRefresher) refreshPosts(ctx context.Context, published bool) error {
	var total int64
	if err := g.db.WithContext(ctx).Model(&models.Post{}).Where("published = ?", published).Count(&total).Error; err != nil {
		return err
	}

	label := "false"
	if published {
		label = "true"
	}
	g.observability.Business.PostsTotal.WithLabelValues(label).Set(float64(total))
	return nil
}

func (g *GaugeRefresher) refreshComments(ctx context.Context) error {
	var total int64
	if err := g.db.WithContext(ctx).Model(&models.Comment{}).Count(&total).Error; err != nil {
		return err
	}

	g.observability.Business.CommentsTotal.Set(float64(total))
	return nil
}
