package service

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/storage"
)

type Service struct {
	OrderService   *OrderService
	ProductService *ProductService
}

func NewService(logger *slog.Logger, repo storage.StorageI) *Service {
	return &Service{
		OrderService:   NewOrderService(logger, repo.OrderRepo(), repo.ProductRepo()),
		ProductService: NewProductService(logger, repo.ProductRepo()),
	}
}
