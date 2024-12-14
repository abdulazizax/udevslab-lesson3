package handler

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/config"
	"github.com/abdulazizax/udevslab-lesson3/internal/service"
)

type Handler struct {
	ProductHandler *ProductHandler
	OrderHandler   *OrderHandler
}

func NewHandler(logger *slog.Logger, service *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		ProductHandler: NewProductHandler(logger, service.ProductService),
		OrderHandler:   NewOrderHandler(logger, service.OrderService),
	}
}
