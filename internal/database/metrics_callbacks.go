package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iruiz/gin-blog-api/internal/metrics"
	"gorm.io/gorm"
)

type dbOp string

const (
	dbCreate dbOp = "create"
	dbQuery  dbOp = "query"
	dbUpdate dbOp = "update"
	dbDelete dbOp = "delete"
)

type callbackStartTimeKey struct{}

type metricsCallbacks struct {
	observability *metrics.Metrics
}

func newMetricsCallbacks(observability *metrics.Metrics) *metricsCallbacks {
	return &metricsCallbacks{observability: observability}
}

func (m *metricsCallbacks) register(db *gorm.DB) error {
	for _, op := range []dbOp{dbCreate, dbQuery, dbUpdate, dbDelete} {
		if err := m.registerOperation(db, op); err != nil {
			return err
		}
	}

	return nil
}

func (m *metricsCallbacks) registerOperation(db *gorm.DB, op dbOp) error {
	switch op {
	case dbCreate:
		if err := db.Callback().Create().Before("gorm:create").Register("metrics:before_create", m.before); err != nil {
			return fmt.Errorf("register before %s callback: %w", op, err)
		}
		if err := db.Callback().Create().After("gorm:create").Register("metrics:after_create", m.after(op)); err != nil {
			return fmt.Errorf("register after %s callback: %w", op, err)
		}
	case dbQuery:
		if err := db.Callback().Query().Before("gorm:query").Register("metrics:before_query", m.before); err != nil {
			return fmt.Errorf("register before %s callback: %w", op, err)
		}
		if err := db.Callback().Query().After("gorm:query").Register("metrics:after_query", m.after(op)); err != nil {
			return fmt.Errorf("register after %s callback: %w", op, err)
		}
	case dbUpdate:
		if err := db.Callback().Update().Before("gorm:update").Register("metrics:before_update", m.before); err != nil {
			return fmt.Errorf("register before %s callback: %w", op, err)
		}
		if err := db.Callback().Update().After("gorm:update").Register("metrics:after_update", m.after(op)); err != nil {
			return fmt.Errorf("register after %s callback: %w", op, err)
		}
	case dbDelete:
		if err := db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", m.before); err != nil {
			return fmt.Errorf("register before %s callback: %w", op, err)
		}
		if err := db.Callback().Delete().After("gorm:delete").Register("metrics:after_delete", m.after(op)); err != nil {
			return fmt.Errorf("register after %s callback: %w", op, err)
		}
	default:
		return fmt.Errorf("unsupported db operation %q", op)
	}

	return nil
}

func (m *metricsCallbacks) before(db *gorm.DB) {
	if db.Statement == nil {
		return
	}

	ctx := db.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}

	db.Statement.Context = context.WithValue(ctx, callbackStartTimeKey{}, time.Now())
}

func (m *metricsCallbacks) after(op dbOp) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if m.observability == nil || db.Statement == nil || db.Statement.Context == nil {
			return
		}

		if start, ok := db.Statement.Context.Value(callbackStartTimeKey{}).(time.Time); ok {
			m.observability.Database.QueryDuration.WithLabelValues(string(op)).Observe(time.Since(start).Seconds())
		}

		if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
			m.observability.Database.ErrorsTotal.WithLabelValues(string(op)).Inc()
		}
	}
}
