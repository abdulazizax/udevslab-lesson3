package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (

	// Products structs

	Product struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Name        string             `bson:"name" json:"name"`
		Description string             `bson:"description" json:"description"`
		Price       float64            `bson:"price" json:"price"`
		Stock       int                `bson:"stock" json:"stock"`
		CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
		UpdatedAt   primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	}

	ProductCreate struct {
		Name        string  `bson:"name" json:"name"`
		Description string  `bson:"description" json:"description"`
		Price       float64 `bson:"price" json:"price"`
		Stock       int     `bson:"stock" json:"stock"`
	}

	ProductUpdate struct {
		Name        string  `bson:"name" json:"name"`
		Description string  `bson:"description" json:"description"`
		Price       float64 `bson:"price" json:"price"`
		Stock       int     `bson:"stock" json:"stock"`
	}

	UpdatedProduct struct {
		Name        string             `bson:"name" json:"name"`
		Description string             `bson:"description" json:"description"`
		Price       float64            `bson:"price" json:"price"`
		Stock       int                `bson:"stock" json:"stock"`
		UpdatedAt   primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	}

	// Orders structs

	Order struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		UserID    primitive.ObjectID `bson:"userId" json:"userId"`
		ProductID primitive.ObjectID `bson:"productId" json:"productId"`
		Quantity  int                `bson:"quantity" json:"quantity"`
		Status    string             `bson:"status" json:"status"`
		Total     float64            `bson:"total" json:"total"`
		CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
		UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	}

	OrderCreate struct {
		UserID    primitive.ObjectID `bson:"userId" json:"userId"`
		ProductID primitive.ObjectID `bson:"productId" json:"productId"`
		Quantity  int                `bson:"quantity" json:"quantity"`
		Status    string             `bson:"status" json:"status"`
	}

	OrderUpdate struct {
		UserID    primitive.ObjectID `bson:"userId" json:"userId"`
		ProductID primitive.ObjectID `bson:"productId,omitempty" json:"productId,omitempty"`
		Quantity  int                `bson:"quantity,omitempty" json:"quantity,omitempty"`
		Status    string             `bson:"status,omitempty" json:"status,omitempty"`
	}

	UpdatedOrder struct {
		UserID    primitive.ObjectID `bson:"userId" json:"userId"`
		ProductID primitive.ObjectID `bson:"productId,omitempty" json:"productId,omitempty"`
		Quantity  int                `bson:"quantity,omitempty" json:"quantity,omitempty"`
		Status    string             `bson:"status,omitempty" json:"status,omitempty"`
		Total     float64            `bson:"total,omitempty" json:"total,omitempty"`
		UpdatedAt primitive.DateTime `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	}

	Report struct {
		TotalProducts int     `json:"totalProducts"`
		TotalOrders   int     `json:"totalOrders"`
		TotalRevenue  float64 `json:"totalRevenue"`
	}

	Error struct {
		Message string `json:"message"`
	}

	Pagination struct {
		Page     int `form:"page" json:"page" binding:"required"`   // Query: ?page=1
		PageSize int `form:"limit" json:"limit" binding:"required"` // Query: ?limit=10
	}

	TopProduct struct {
		ProductID    string  `json:"product_id" bson:"_id"`
		Name         string  `json:"name" bson:"name"`
		TotalSold    int     `json:"total_sold" bson:"total_sold"`
		TotalRevenue float64 `json:"total_revenue" bson:"total_revenue"`
	}

	OrderWithUser struct {
		OrderID   string    `json:"order_id" bson:"_id"`
		UserID    string    `json:"user_id" bson:"user_id"`
		UserName  string    `json:"user_name" bson:"user_name"`
		ProductID string    `json:"product_id" bson:"product_id"`
		Quantity  int       `json:"quantity" bson:"quantity"`
		Status    string    `json:"status" bson:"status"`
		Total     float64   `json:"total" bson:"total"`
		CreatedAt time.Time `json:"created_at" bson:"created_at"`
	}

	OrderAggregate struct {
		Date         string  `json:"date" bson:"_id"`
		TotalOrders  int     `json:"total_orders" bson:"total_orders"`
		TotalRevenue float64 `json:"total_revenue" bson:"total_revenue"`
	}
)
