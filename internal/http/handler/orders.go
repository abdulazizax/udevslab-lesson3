package handler

import (
	"log/slog"
	"net/http"

	"github.com/abdulazizax/udevslab-lesson3/internal/models"
	"github.com/abdulazizax/udevslab-lesson3/internal/service"
	"github.com/gin-gonic/gin"
)

// OrderHandler handles the HTTP requests for orders
type OrderHandler struct {
	logger       *slog.Logger
	orderService *service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(logger *slog.Logger, orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		logger:       logger,
		orderService: orderService,
	}
}

// CreateOrder creates a new order
// @Summary Create a new order
// @Description Create a new order and return the created order's ID
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param order body models.Order true "Order information"
// @Success 201 {string} string "Order ID"
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /orders [post]
func (s *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		s.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid order data"})
		return
	}

	orderID, err := s.orderService.CreateOrder(c, &order)
	if err != nil {
		s.logger.Error("failed to create order", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, orderID)
}

// ListOrders lists all orders
// @Summary List all orders
// @Description Retrieve all orders from the system
// @Tags Orders
// @Produce  json
// @Success 200 {array} models.Order "List of orders"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /orders [get]
func (s *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := s.orderService.ListOrders(c)
	if err != nil {
		s.logger.Error("failed to list orders", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder fetches an order by ID
// @Summary Get an order by ID
// @Description Fetch a single order by its ID
// @Tags Orders
// @Produce  json
// @Param order_id path string true "Order ID"
// @Success 200 {object} models.Order "Order details"
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 404 {object} models.Error "Order Not Found"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /orders/{order_id} [get]
func (s *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("order_id")

	order, err := s.orderService.GetOrderByID(c, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, models.Error{Message: "Order not found"})
		} else {
			s.logger.Error("failed to get order", "error", err)
			c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to fetch order"})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrder updates an existing order
// @Summary Update an order
// @Description Update the details of an existing order
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param updates body models.OrderUpdate true "Order fields to update"
// @Success 200 {string} string "Order updated successfully"
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 404 {object} models.Error "Order Not Found"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /orders/{order_id} [put]
func (s *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	var updates models.OrderUpdate
	if err := c.ShouldBindJSON(&updates); err != nil {
		s.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid data for update"})
		return
	}

	err := s.orderService.UpdateOrder(c, orderID, &updates)
	if err != nil {
		if err.Error() == "no order found to update" {
			c.JSON(http.StatusNotFound, models.Error{Message: "Order not found"})
		} else {
			s.logger.Error("failed to update order", "error", err)
			c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to update order"})
		}
		return
	}

	c.JSON(http.StatusOK, "Order updated successfully")
}

// DeleteOrder deletes an order
// @Summary Delete an order
// @Description Delete an order by its ID
// @Tags Orders
// @Produce  json
// @Param order_id path string true "Order ID"
// @Success 200 {string} string "Order deleted successfully"
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 404 {object} models.Error "Order Not Found"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /orders/{order_id} [delete]
func (s *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("order_id")

	err := s.orderService.DeleteOrder(c, orderID)
	if err != nil {
		if err.Error() == "no order found to delete" {
			c.JSON(http.StatusNotFound, models.Error{Message: "Order not found"})
		} else {
			s.logger.Error("failed to delete order", "error", err)
			c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to delete order"})
		}
		return
	}

	c.JSON(http.StatusOK, "Order deleted successfully")
}
