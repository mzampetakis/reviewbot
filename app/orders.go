package app

import (
	"context"
	"time"
)

type OrderStatus string

const (
	OrderStatusPlaced    OrderStatus = "placed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusSending   OrderStatus = "sending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusReviewed  OrderStatus = "reviewed"
)

// Customer represents a customer entity.
type Customer struct {
	UUID             string    `json:"uuid"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	PhoneNumber      string    `json:"phone_number"`
	RegistrationDate time.Time `json:"registration_date"`
}

// Order represents an order entity at the Database.
type Order struct {
	UUID       string      `json:"uuid"`
	Customer   Customer    `json:"customer"`
	Status     OrderStatus `json:"status"`
	PlacedDate time.Time   `json:"placed_date"`
}

// OrderProduct represents an order products entity.
type OrderProduct struct {
	UUID        string  `json:"uuid"`
	OrderUUID   string  `json:"order_uuid"`
	ProductUUID string  `json:"product_uuid"`
	Items       int     `json:"items"`
	Product     Product `json:"product"`
}

// Product represents a product entity.
type Product struct {
	UUID               string    `json:"uuid"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Image              string    `json:"items"`
	AvailabilityStatus string    `json:"availability_status"`
	AvailableItems     int       `json:"available_items"`
	CreatedAt          time.Time `json:"createdAt"`
	Manufacturer       string    `json:"manufacturer"`
	Vehicle            string    `json:"vehicle"`
	ID                 string    `json:"id"`
}

// OrdersRepository should be implemented to get access to the data store.
type OrdersRepository interface {
	GetOrderByUUID(ctx context.Context, uuid string) (*Order, error)
	UpdateOrderStatusByOrderUUID(ctx context.Context, uuid string, status string) error
	GetOrderProductsByOrderUUID(ctx context.Context, uuid string) ([]OrderProduct, error)
	AddOrderProductReviewByOrderProductUUID(ctx context.Context, orderProductUUID string,
		score int64) error
	AddProduct(ctx context.Context, product Product) error
}
