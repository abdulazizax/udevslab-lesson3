package handler

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	logger       *slog.Logger
	orderService *service.OrderService
}

func NewOrderHandler(logger *slog.Logger, orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		logger:       logger,
		orderService: orderService,
	}
}

func (s *OrderHandler) CreateOrder(c *gin.Context) {}
func (s *OrderHandler) ListOrders(c *gin.Context)  {}
func (s *OrderHandler) GetOrder(c *gin.Context)    {}
func (s *OrderHandler) UpdateOrder(c *gin.Context) {}
func (s *OrderHandler) DeleteOrder(c *gin.Context) {}
