package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/abdulazizax/udevslab-lesson3/internal/models"
	"github.com/abdulazizax/udevslab-lesson3/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductHandler struct holds the logger and product service.
type ProductHandler struct {
	logger         *slog.Logger
	productService *service.ProductService
}

// NewProductHandler creates a new instance of ProductHandler.
func NewProductHandler(logger *slog.Logger, productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		logger:         logger,
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Creates a new product in the database
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.ProductCreate true "Product information"
// @Success 201 {object} gin.H "Product created successfully"// @Failure 400 {object} models.Error "Bad request"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products [post]
func (s *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.ProductCreate
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request body"})
		return
	}

	productID, err := s.productService.CreateProduct(c, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	id, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ListProducts godoc
// @Summary List all products
// @Description Retrieve a list of all products in the database
// @Tags products
// @Produce json
// @Param page query int true "Page number"
// @Param limit query int true "Items per page"
// @Success 200 {array} models.Product "List of products"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products [get]
func (s *ProductHandler) ListProducts(c *gin.Context) {
	var pagination models.Pagination
	// Bind query parameters to the Pagination struct
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid query parameters"})
		return
	}

	// Call the service layer to get paginated products
	products, err := s.productService.ListProducts(c, &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	// Return paginated products as JSON
	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Retrieve a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product "Product found"
// @Failure 404 {object} models.Error "Product not found"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products/{id} [get]
func (s *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	product, err := s.productService.GetProductByID(c, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Update product details by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.ProductUpdate true "Updated product details"
// @Success 200 {object} gin.H "Product updated successfully"
// @Failure 400 {object} models.Error "Bad request"
// @Failure 404 {object} models.Error "Product not found"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products/{id} [put]
func (s *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	var product models.ProductUpdate
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request body"})
		return
	}

	err := s.productService.UpdateProduct(c, productID, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Poduct updated successfully"})
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a product from the database by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 {object} gin.H "Product deleted successfully"
// @Failure 404 {object} models.Error "Product not found"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products/{id} [delete]
func (s *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	err := s.productService.DeleteProduct(c, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Product deleted successfully"})
}

// SearchProductsByName godoc
// @Summary Search products by name
// @Description Search products by partial name with pagination
// @Tags products
// @Produce json
// @Param name query string true "Search keyword"
// @Param page query int true "Page number"
// @Param limit query int true "Items per page"
// @Success 200 {array} models.Product "List of products"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products/search [get]
func (s *ProductHandler) SearchProductsByName(c *gin.Context) {
	var pagination models.Pagination
	searchQuery := c.Query("name") // Get the search keyword from query
	if searchQuery == "" {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Search query is required"})
		return
	}

	// Bind pagination query parameters
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid pagination parameters"})
		return
	}

	// Call the service layer to get search results with pagination
	products, err := s.productService.SearchProductsByName(c, searchQuery, &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	// Return search results
	c.JSON(http.StatusOK, products)
}

// ExactSearchProductsByPrice godoc
// @Summary Search products by exact price
// @Description Retrieve products based on an exact price match with pagination
// @Tags products
// @Produce json
// @Param price query float64 true "Price to search for"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of products per page" default(10)
// @Success 200 {array} models.Product "List of products"
// @Failure 400 {object} models.Error "Invalid request parameters"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products/search/price [get]
func (s *ProductHandler) ExactSearchProductsByPrice(c *gin.Context) {
	// Get price from query parameter
	priceParam := c.DefaultQuery("price", "0") // Default is 0 if not provided
	price, err := strconv.ParseFloat(priceParam, 64)
	if err != nil || price <= 0 {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid price parameter"})
		return
	}

	// Get pagination parameters
	var pagination models.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid pagination parameters"})
		return
	}

	// Call the service method for fetching products by exact price
	products, err := s.productService.ExactSearchProductsByPrice(c, price, &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	// Return the products in JSON response
	c.JSON(http.StatusOK, products)
}

// SearchProductsByPriceRangeInc godoc
// @Summary Search products by price range
// @Description Retrieve products based on a price range with pagination
// @Tags products
// @Produce json
// @Param order query int8 true "Order (-1: decreasing, 1: increasing)"
// @Param min_price query float64 true "Minimum price"
// @Param max_price query float64 true "Maximum price"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of products per page" default(10)
// @Success 200 {array} models.Product "List of products"
// @Failure 400 {object} models.Error "Invalid request parameters"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /products/search/price-range [get]
func (s *ProductHandler) SearchProductsByPriceRange(c *gin.Context) {
	// Get price range from query parameters
	orderStr := c.DefaultQuery("order", "1")          // Default is 1 if not provided
	minPriceParam := c.DefaultQuery("min_price", "0") // Default is 0 if not provided
	maxPriceParam := c.DefaultQuery("max_price", "0") // Default is 0 if not provided

	order, err := strconv.ParseInt(orderStr, 10, 8)
	if err != nil || order < -1 || order > 1 {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid order parameter. Must be -1 or 1"})
		return
	}

	minPrice, err := strconv.ParseFloat(minPriceParam, 64)
	if err != nil || minPrice <= 0 {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid minimum price parameter"})
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceParam, 64)
	if err != nil || maxPrice <= 0 {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid maximum price parameter"})
		return
	}

	// Get pagination parameters
	var pagination models.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid pagination parameters"})
		return
	}

	// Call the service method for fetching products by price range
	products, err := s.productService.SearchProductsByPriceRange(c, int8(order), minPrice, maxPrice, &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	// Return the products in JSON response
	c.JSON(http.StatusOK, products)
}
