package orders

import (
	"context"
	"golang.org/x/exp/slog"
	"reviewbot/app"
)

// Service wraps the user repository.
type Service struct {
	repo   app.OrdersRepository
	logger *slog.Logger
}

// NewService returns a new Service.
func NewService(repo app.OrdersRepository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// OrderByUUID gets an order by its UUID.
func (s *Service) OrderByUUID(ctx context.Context, orderUUID string) (*app.Order, error) {
	return s.repo.GetOrderByUUID(ctx, orderUUID)
}

// UpdateOrderStatusByUUID updates an order status by its UUID.
func (s *Service) UpdateOrderStatusByUUID(ctx context.Context, orderUUID string,
	orderStatus app.OrderStatus) error {
	return s.repo.UpdateOrderStatusByOrderUUID(ctx, orderUUID, string(orderStatus))
}

// OrderProductsByOrderUUID gets an order by its UUID.
func (s *Service) OrderProductsByOrderUUID(ctx context.Context, orderUUID string) ([]app.OrderProduct, error) {
	return s.repo.GetOrderProductsByOrderUUID(ctx, orderUUID)
}
