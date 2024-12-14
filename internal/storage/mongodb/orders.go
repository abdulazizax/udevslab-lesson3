package mongodb

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderSrorage struct {
	db     *mongo.Collection
	logger *slog.Logger
	cfg    *config.Config
}

func NewOrderSrorage(db *mongo.Database, logger *slog.Logger, cfg *config.Config) *OrderSrorage {
	return &OrderSrorage{
		db:     db.Collection("Orders"),
		logger: logger,
		cfg:    cfg,
	}
}
