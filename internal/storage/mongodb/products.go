package mongodb

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductSrorage struct {
	db     *mongo.Collection
	logger *slog.Logger
	cfg    *config.Config
}

func NewProductSrorage(db *mongo.Database, logger *slog.Logger, cfg *config.Config) *ProductSrorage {
	return &ProductSrorage{
		db:     db.Collection("Products"),
		logger: logger,
		cfg:    cfg,
	}
}
