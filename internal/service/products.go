package service

import (
	"context"
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/models"
	"github.com/abdulazizax/udevslab-lesson3/internal/repos"
)

type ProductService struct {
	logger      *slog.Logger
	productRepo repos.ProductRepo
}

func NewProductService(logger *slog.Logger, productRepo repos.ProductRepo) *ProductService {
	return &ProductService{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.ProductCreate) (string, error) {
	return s.productRepo.CreateProduct(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, productID string) (*models.Product, error) {
	return s.productRepo.GetProductByID(ctx, productID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID string, updates *models.ProductUpdate) error {
	return s.productRepo.UpdateProduct(ctx, productID, updates)
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID string) error {
	return s.productRepo.DeleteProduct(ctx, productID)
}

func (s *ProductService) ListProducts(ctx context.Context, pagination *models.Pagination) ([]models.Product, error) {
	return s.productRepo.ListProducts(ctx, pagination)
}

func (s *ProductService) SearchProductsByName(ctx context.Context, name string, pagination *models.Pagination) ([]models.Product, error) {
	return s.productRepo.SearchProductsByName(ctx, name, pagination)
}

func (s *ProductService) ExactSearchProductsByPrice(ctx context.Context, price float64, pagination *models.Pagination) ([]models.Product, error) {
	return s.productRepo.ExactSearchProductsByPrice(ctx, price, pagination)
}

func (s *ProductService) SearchProductsByPriceRange(ctx context.Context, order int8, minPrice, maxPrice float64, pagination *models.Pagination) ([]models.Product, error) {
	return s.productRepo.SearchProductsByPriceRange(ctx, order, minPrice, maxPrice, pagination)
}
