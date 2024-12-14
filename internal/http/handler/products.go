package handler

import (
	"log/slog"

	"github.com/abdulazizax/udevslab-lesson3/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	logger         *slog.Logger
	productService *service.ProductService
}

func NewProductHandler(logger *slog.Logger, productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		logger:         logger,
		productService: productService,
	}
}

func (s *ProductHandler) CreateProduct(c *gin.Context) {}
func (s *ProductHandler) ListProducts(c *gin.Context)  {}
func (s *ProductHandler) GetProduct(c *gin.Context)    {}
func (s *ProductHandler) UpdateProduct(c *gin.Context) {}
func (s *ProductHandler) DeleteProduct(c *gin.Context) {}
