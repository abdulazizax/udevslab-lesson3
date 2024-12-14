package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/abdulazizax/udevslab-lesson3/internal/config"
	"github.com/abdulazizax/udevslab-lesson3/internal/models"
	"github.com/abdulazizax/udevslab-lesson3/internal/repos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductStorage struct {
	db     *mongo.Collection
	logger *slog.Logger
	cfg    *config.Config
}

func NewProductStorage(db *mongo.Database, logger *slog.Logger, cfg *config.Config) repos.ProductRepo {
	return &ProductStorage{
		db:     db.Collection("Products"),
		logger: logger,
		cfg:    cfg,
	}
}

// CreateProduct creates a new product in the database
func (p *ProductStorage) CreateProduct(ctx context.Context, product *models.ProductCreate) (string, error) {
	p.logger.Info("starting product creation", "name", product.Name)

	created_at := time.Now()

	// Insert the product into MongoDB
	var newProduct = models.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   primitive.NewDateTimeFromTime(created_at),
		UpdatedAt:   primitive.NewDateTimeFromTime(created_at),
	}
	result, err := p.db.InsertOne(ctx, newProduct)
	if err != nil {
		p.logger.Error("failed to insert product", "error", err)
		return "", fmt.Errorf("failed to insert product: %w", err)
	}

	// Get the inserted ID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		p.logger.Error("failed to convert inserted ID to ObjectID")
		return "", errors.New("failed to convert inserted ID to ObjectID")
	}

	p.logger.Info("product creation successful", "productID", insertedID.Hex())
	return insertedID.Hex(), nil
}

// GetProductByID fetches a product by its ID
func (p *ProductStorage) GetProductByID(ctx context.Context, productID string) (*models.Product, error) {
	p.logger.Info("fetching product by ID", "productID", productID)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		p.logger.Error("invalid product ID format", "error", err)
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}

	// Find the product in MongoDB
	var product models.Product
	err = p.db.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			p.logger.Warn("product not found", "productID", productID)
			return nil, errors.New("product not found")
		}
		p.logger.Error("failed to fetch product from database", "error", err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	p.logger.Info("successfully fetched product", "productID", productID)
	return &product, nil
}

// UpdateProduct updates a product in the database
func (p *ProductStorage) UpdateProduct(ctx context.Context, productID string, updates *models.ProductUpdate) error {
	p.logger.Info("updating product", "productID", productID)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		p.logger.Error("invalid product ID format", "error", err)
		return fmt.Errorf("invalid product ID format: %w", err)
	}

	var newProduct = models.UpdatedProduct{
		Name:        updates.Name,
		Description: updates.Description,
		Price:       updates.Price,
		Stock:       updates.Stock,
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	// Update the product in MongoDB
	result, err := p.db.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": newProduct})
	if err != nil {
		p.logger.Error("failed to update product", "error", err)
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.MatchedCount == 0 {
		p.logger.Warn("no product found to update", "productID", productID)
		return errors.New("no product found to update")
	}

	p.logger.Info("product updated successfully", "productID", productID, "updatedFields", updates)
	return nil
}

// DeleteProduct deletes a product by its ID
func (p *ProductStorage) DeleteProduct(ctx context.Context, productID string) error {
	p.logger.Info("deleting product", "productID", productID)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		p.logger.Error("invalid product ID format", "error", err)
		return fmt.Errorf("invalid product ID format: %w", err)
	}

	// Delete the product in MongoDB
	result, err := p.db.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		p.logger.Error("failed to delete product", "error", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.DeletedCount == 0 {
		p.logger.Warn("no product found to delete", "productID", productID)
		return errors.New("no product found to delete")
	}

	p.logger.Info("product deleted successfully", "productID", productID)
	return nil
}

// ListProducts fetches a paginated list of products from the database
func (p *ProductStorage) ListProducts(ctx context.Context, pagination *models.Pagination) ([]models.Product, error) {
	p.logger.Info("fetching list of products", "page", pagination.Page, "pageSize", pagination.Page)

	// Hisoblash: nechta yozuvni o'tkazib yuborish kerak
	skip := (pagination.Page - 1) * pagination.PageSize

	// MongoDB query uchun options
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.PageSize))

	// Find products with pagination
	cursor, err := p.db.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		p.logger.Error("failed to fetch products from database", "error", err)
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer cursor.Close(ctx)

	// Ma'lumotlarni massivga tushirish
	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		p.logger.Error("failed to decode products", "error", err)
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	p.logger.Info("successfully fetched products", "productCount", len(products))
	return products, nil
}

// SearchProducts searches for products by name with pagination
func (p *ProductStorage) SearchProductsByName(ctx context.Context, name string, pagination *models.Pagination) ([]models.Product, error) {
	// Define the search criteria using $regex for partial matching
	filter := bson.M{
		"name": bson.M{
			"$regex":   name, // Partial match
			"$options": "i",  // Case-insensitive
		},
	}

	// Pagination settings
	skip := (pagination.Page - 1) * pagination.PageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.PageSize))

	// Execute the query
	cursor, err := p.db.Find(ctx, filter, opts)
	if err != nil {
		p.logger.Error("failed to search products from database", "error", err)
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode the results
	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		p.logger.Error("failed to decode search results", "error", err)
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	p.logger.Info("successfully fetched search results", "productCount", len(products))
	return products, nil
}

// ExactSearchProducts searches for products by price with pagination
func (s *ProductStorage) ExactSearchProductsByPrice(ctx context.Context, price float64, pagination *models.Pagination) ([]models.Product, error) {
	// Create filter for exact price match
	filter := bson.M{"price": price}

	// Pagination and sorting logic
	skip := (pagination.Page - 1) * pagination.PageSize
	options := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.PageSize)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}}) // Sort by createdAt, descending

	cursor, err := s.db.Find(ctx, filter, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// SearchProductsByPriceRangeInc searches for products by price with pagination
func (s *ProductStorage) SearchProductsByPriceRange(ctx context.Context, order int8, minPrice, maxPrice float64, pagination *models.Pagination) ([]models.Product, error) {
	// Create filter for price range match
	filter := bson.M{
		"price": bson.M{
			"$gte": minPrice, // Greater than or equal to minPrice
			"$lte": maxPrice, // Less than or equal to maxPrice
		},
	}

	// Pagination and sorting logic
	skip := (pagination.Page - 1) * pagination.PageSize
	options := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.PageSize)).
		SetSort(bson.D{{Key: "price", Value: int(order)}}) // Sort by price in ascending order

	cursor, err := s.db.Find(ctx, filter, options)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}
