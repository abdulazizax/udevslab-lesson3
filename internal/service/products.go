package service

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/repos"
)

type ProductService struct {
	logger      *slog.Logger
	productRepo *repos.ProductRepo
}

func NewProductService(logger *slog.Logger, productRepo *repos.ProductRepo) *ProductService {
	return &ProductService{
		logger:      logger,
		productRepo: productRepo,
	}
}
