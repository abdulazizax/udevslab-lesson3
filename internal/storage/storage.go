package storage

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/config"
	"github.com/abdulazizax/udevslab-lesson3/internal/repos"
	"github.com/abdulazizax/udevslab-lesson3/internal/storage/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageI interface {
	ProductRepo() repos.ProductRepo
	OrderRepo() repos.OrderRepo
}

type Storage struct {
	productRepo repos.ProductRepo
	orderRepo   repos.OrderRepo
}

func New(db *mongo.Database, cfg *config.Config, logger *slog.Logger) StorageI {
	return &Storage{
		productRepo: mongodb.NewProductSrorage(db, logger, cfg),
		orderRepo:   mongodb.NewOrderSrorage(db, logger, cfg),
	}
}

func (s *Storage) ProductRepo() repos.ProductRepo {
	return s.productRepo
}

func (s *Storage) OrderRepo() repos.OrderRepo {
	return s.orderRepo
}
