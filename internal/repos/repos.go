package repos

import (
	"context"
	"time"

	"github.com/abdulazizax/udevslab-lesson3/internal/models"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, total float64, order *models.OrderCreate) (string, error)
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	UpdateOrder(ctx context.Context, total float64, orderID string, updates *models.OrderUpdate) (string, error)
	DeleteOrder(ctx context.Context, orderID string) error
	ListOrders(ctx context.Context, pagination *models.Pagination) ([]models.Order, error)
	ListOrdersByDateRange(ctx context.Context, order int8, pagination *models.Pagination, startDate, endDate time.Time) ([]models.Order, error)
	ListOrdersWithAggregates(ctx context.Context, startDate, endDate time.Time, pagination *models.Pagination, order int8) ([]models.OrderAggregate, error)
	ListOrdersByCustomer(ctx context.Context, customerID string) ([]models.Order, error)
}

type ProductRepo interface {
	CreateProduct(ctx context.Context, product *models.ProductCreate) (string, error)
	GetProductByID(ctx context.Context, productID string) (*models.Product, error)
	UpdateProduct(ctx context.Context, productID string, updates *models.ProductUpdate) error
	DeleteProduct(ctx context.Context, productID string) error
	ListProducts(ctx context.Context, pagination *models.Pagination) ([]models.Product, error)
	SearchProductsByName(ctx context.Context, name string, pagination *models.Pagination) ([]models.Product, error)
	ExactSearchProductsByPrice(ctx context.Context, price float64, pagination *models.Pagination) ([]models.Product, error)
	SearchProductsByPriceRange(ctx context.Context, order int8, minPrice, maxPrice float64, pagination *models.Pagination) ([]models.Product, error)
	TopSellingProducts(ctx context.Context) ([]models.ProductSales, error)
}
