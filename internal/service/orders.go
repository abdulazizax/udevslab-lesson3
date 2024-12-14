package service

import (
	"context"
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/models"
	"github.com/abdulazizax/udevslab-lesson3/internal/repos"
)

type OrderService struct {
	logger    *slog.Logger
	orderRepo repos.OrderRepo
}

func NewOrderService(logger *slog.Logger, orderRepo repos.OrderRepo) *OrderService {
	return &OrderService{
		logger:    logger,
		orderRepo: orderRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) (string, error) {
	return s.orderRepo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(ctx, orderID)
}

func (s *OrderService) UpdateOrder(ctx context.Context, orderID string, updates *models.OrderUpdate) error {
	return s.orderRepo.UpdateOrder(ctx, orderID, updates)
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderID string) error {
	return s.orderRepo.DeleteOrder(ctx, orderID)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]models.Order, error) {
	return s.orderRepo.ListOrders(ctx)
}
