package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/abdulazizax/udevslab-lesson3/internal/models"
	"github.com/abdulazizax/udevslab-lesson3/internal/repos"
)

type OrderService struct {
	logger      *slog.Logger
	orderRepo   repos.OrderRepo
	productRepo repos.ProductRepo
}

func NewOrderService(logger *slog.Logger, orderRepo repos.OrderRepo, productRepo repos.ProductRepo) *OrderService {
	return &OrderService{
		logger:      logger,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.OrderCreate) (string, error) {
	product, err := s.productRepo.GetProductByID(ctx, order.ProductID.Hex())
	if err != nil {
		return "Product is not exists", err
	} else {
		return s.orderRepo.CreateOrder(ctx, product.Price*float64(order.Quantity), order)
	}
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(ctx, orderID)
}

func (s *OrderService) UpdateOrder(ctx context.Context, orderID string, updates *models.OrderUpdate) (string, error) {
	product, err := s.productRepo.GetProductByID(ctx, updates.ProductID.Hex())
	if err != nil {
		return "Product is not exists", err
	} else {
		return s.orderRepo.UpdateOrder(ctx, product.Price*float64(updates.Quantity), orderID, updates)
	}
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderID string) error {
	return s.orderRepo.DeleteOrder(ctx, orderID)
}

func (s *OrderService) ListOrders(ctx context.Context, pagination *models.Pagination) ([]models.Order, error) {
	return s.orderRepo.ListOrders(ctx, pagination)
}

func (s *OrderService) ListOrdersByDateRange(ctx context.Context, order int8, pagination *models.Pagination, startDate, endDate time.Time) ([]models.Order, error) {
	return s.orderRepo.ListOrdersByDateRange(ctx, order, pagination, startDate, endDate)
}