package service

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/repos"
)

type OrderService struct {
	logger    *slog.Logger
	orderRepo *repos.OrderRepo
}

func NewOrderService(logger *slog.Logger, orderRepo *repos.OrderRepo) *OrderService {
	return &OrderService{
		logger:    logger,
		orderRepo: orderRepo,
	}
}
