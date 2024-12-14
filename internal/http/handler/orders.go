package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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
// @Param order body models.OrderCreate true "Order information"
// @Success 201 {object} gin.H "Order ID"
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /orders [post]
func (s *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.OrderCreate
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

	c.JSON(http.StatusCreated, gin.H{"id": orderID})
}

// ListOrders godoc
// @Summary List all orders
// @Description Retrieve a list of all orders in the database with pagination
// @Tags Orders
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {array} models.Order "List of orders"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /orders [get]
func (o *OrderHandler) ListOrders(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number."})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size."})
		return
	}

	// Create pagination parameters
	var pagination = &models.Pagination{
		Page:     pageInt,
		PageSize: pageSizeInt,
	}

	orders, err := o.orderService.ListOrders(c, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
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
// @Success 200 {object} gin.H "Order updated successfully"
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

	res, err := s.orderService.UpdateOrder(c, orderID, &updates)
	if res != "" {
		c.JSON(http.StatusNotFound, models.Error{Message: "Order updated successfully"})
	}
	if err != nil {
		if err.Error() == "no order found to update" {
			c.JSON(http.StatusNotFound, models.Error{Message: "Order not found"})
		} else {
			s.logger.Error("failed to update order", "error", err)
			c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to update order"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
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

// ListOrders godoc
// @Summary List orders within a specific date range
// @Description Retrieve a paginated list of orders filtered by a specific date range and sorted by the creation date in ascending or descending order.
// @Tags orders
// @Produce json
// @Param order query int true "Sort order (-1: descending, 1: ascending)" default(1)
// @Param page query int false "Page number (1-based index)" default(1)
// @Param page_size query int false "Number of orders per page" default(10)
// @Param start_date query string true "Start date in format (YYYY-MM-DD)" default(2000-01-01)
// @Param end_date query string true "End date in format (YYYY-MM-DD)" default(2026-01-01)
// @Success 200 {array} models.Order "Paginated list of orders"
// @Failure 400 {object} models.Error "Bad request (invalid parameters)"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /orders/range [get]
func (o *OrderHandler) ListOrdersByDateRange(c *gin.Context) {
	// Get start and end dates from query parameters
	orderStr := c.DefaultQuery("order", "1") // Default is 1 if not provided
	startDateStr := c.DefaultQuery("start_date", "")
	endDateStr := c.DefaultQuery("end_date", "")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number."})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size."})
		return
	}

	// Create pagination parameters
	var pagination = &models.Pagination{
		Page:     pageInt,
		PageSize: pageSizeInt,
	}

	order, err := strconv.ParseInt(orderStr, 10, 8)
	if err != nil || order < -1 || order > 1 {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid order parameter. Must be -1 or 1"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid end date format"})
		return
	}

	// Fetch orders with pagination, date range filter, and sorting by date
	orders, err := o.orderService.ListOrdersByDateRange(c, int8(order), pagination, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	fmt.Println(order)

	c.JSON(http.StatusOK, orders)
}
