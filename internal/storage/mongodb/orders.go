package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdulazizax/udevslab-lesson3/internal/config"
	"github.com/abdulazizax/udevslab-lesson3/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderStorage struct {
	db     *mongo.Collection
	logger *slog.Logger
	cfg    *config.Config
}

func NewOrderStorage(db *mongo.Database, logger *slog.Logger, cfg *config.Config) *OrderStorage {
	return &OrderStorage{
		db:     db.Collection("Orders"),
		logger: logger,
		cfg:    cfg,
	}
}

// CreateOrder creates a new order in the database
func (o *OrderStorage) CreateOrder(ctx context.Context, total float64, order *models.OrderCreate) (string, error) {
	o.logger.Info("starting order creation", "ProductID", order.ProductID)

	created_at := time.Now()

	// Insert the order into MongoDB
	var newOrder = models.Order{
		UserID:    order.UserID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		Status:    order.Status,
		Total:     total,
		CreatedAt: primitive.NewDateTimeFromTime(created_at),
		UpdatedAt: primitive.NewDateTimeFromTime(created_at),
	}

	result, err := o.db.InsertOne(ctx, newOrder)
	if err != nil {
		o.logger.Error("failed to insert order", "error", err)
		return "", fmt.Errorf("failed to insert order: %w", err)
	}

	// Get the inserted ID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		o.logger.Error("failed to convert inserted ID to ObjectID")
		return "", errors.New("failed to convert inserted ID to ObjectID")
	}

	o.logger.Info("order creation successful", "orderID", insertedID.Hex())
	return insertedID.Hex(), nil
}

// GetOrderByID fetches an order by its ID
func (o *OrderStorage) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	o.logger.Info("fetching order by ID", "orderID", orderID)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		o.logger.Error("invalid order ID format", "error", err)
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	// Find the order in MongoDB
	var order models.Order
	err = o.db.FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			o.logger.Warn("order not found", "orderID", orderID)
			return nil, errors.New("order not found")
		}
		o.logger.Error("failed to fetch order from database", "error", err)
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	o.logger.Info("successfully fetched order", "orderID", orderID)
	return &order, nil
}

// UpdateOrder updates an order in the database
func (o *OrderStorage) UpdateOrder(ctx context.Context, total float64, orderID string, updates *models.OrderUpdate) (string, error) {
	o.logger.Info("updating order", "orderID", orderID)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		o.logger.Error("invalid order ID format", "error", err)
		return "", fmt.Errorf("invalid order ID format: %w", err)
	}

	var newOrder = models.UpdatedOrder{
		UserID:    updates.UserID,
		ProductID: updates.ProductID,
		Quantity:  updates.Quantity,
		Status:    updates.Status,
		Total:     total,
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now())}

	// Update the order in MongoDB
	result, err := o.db.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": newOrder})
	if err != nil {
		o.logger.Error("failed to update order", "error", err)
		return "", fmt.Errorf("failed to update order: %w", err)
	}

	if result.MatchedCount == 0 {
		o.logger.Warn("no order found to update", "orderID", orderID)
		return "", errors.New("no order found to update")
	}

	o.logger.Info("order updated successfully", "orderID", orderID, "updatedFields", updates)
	return "", nil
}

// DeleteOrder deletes an order by its ID
func (o *OrderStorage) DeleteOrder(ctx context.Context, orderID string) error {
	o.logger.Info("deleting order", "orderID", orderID)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		o.logger.Error("invalid order ID format", "error", err)
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	// Delete the order in MongoDB
	result, err := o.db.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		o.logger.Error("failed to delete order", "error", err)
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.DeletedCount == 0 {
		o.logger.Warn("no order found to delete", "orderID", orderID)
		return errors.New("no order found to delete")
	}

	o.logger.Info("order deleted successfully", "orderID", orderID)
	return nil
}

// ListOrders fetches all orders from the database with pagination
func (o *OrderStorage) ListOrders(ctx context.Context, pagination *models.Pagination) ([]models.Order, error) {
	o.logger.Info("fetching all orders with pagination")

	// Find all orders in MongoDB
	skip := (pagination.Page - 1) * pagination.PageSize
	options := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.PageSize))

	cursor, err := o.db.Find(ctx, bson.M{}, options)
	if err != nil {
		o.logger.Error("failed to fetch orders from database", "error", err)
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each order
	var orders []models.Order
	for cursor.Next(ctx) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			o.logger.Error("failed to decode order", "error", err)
			return nil, fmt.Errorf("failed to decode order: %w", err)
		}
		orders = append(orders, order)
	}

	// Check for errors that occurred during iteration
	if err := cursor.Err(); err != nil {
		o.logger.Error("cursor error", "error", err)
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	o.logger.Info("successfully fetched orders with pagination", "orderCount", len(orders))
	return orders, nil
}

// ListOrdersByDateRange - Get orders by date range (from startDate to endDate) with pagination
func (s *OrderStorage) ListOrdersByDateRange(ctx context.Context, order int8, pagination *models.Pagination, startDate, endDate time.Time) ([]models.Order, error) {
	var orders []models.Order

	// Calculate skip and limit for pagination
	skip := (pagination.Page - 1) * pagination.PageSize
	limit := pagination.PageSize

	// Aggregation pipeline to filter orders by date range and apply pagination
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "createdAt", Value: bson.D{
				{Key: "$gte", Value: startDate}, // greater than or equal to startDate
				{Key: "$lte", Value: endDate},   // less than or equal to endDate
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "createdAt", Value: int(order)}, // Sort by createdAt in ascending or descending order
		}}},
		bson.D{{Key: "$skip", Value: skip}},   // Skip the number of records for pagination
		bson.D{{Key: "$limit", Value: limit}}, // Limit the number of records per page
	}

	cursor, err := s.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // Close the cursor once done

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errors.New("no orders found in the given date range")
	}

	return orders, nil
}
