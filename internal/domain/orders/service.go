package orders

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
	"net/http"
	"reviewbot/app"
	"reviewbot/pkg/responsegenerator"
	"reviewbot/pkg/sentimentanalyzer"
	"time"
)

// Service wraps the user repository.
type Service struct {
	repo              app.OrdersRepository
	responseGenerator responsegenerator.ResponseGenerator
	sentimentAnalyzer sentimentanalyzer.SentimentAnalyze
	logger            *slog.Logger
}

// NewService returns a new Service.
func NewService(repo app.OrdersRepository, responseGenerator responsegenerator.ResponseGenerator,
	sentimentAnalyzer sentimentanalyzer.SentimentAnalyze,
	logger *slog.Logger) *Service {
	return &Service{repo: repo, responseGenerator: responseGenerator, sentimentAnalyzer: sentimentAnalyzer,
		logger: logger}
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

type RemoteProduct struct {
	CreatedAt    time.Time `json:"createdAt"`
	ProductName  string    `json:"productName"`
	Manufacturer string    `json:"manufacturer"`
	Vehicle      string    `json:"vehicle"`
	ImageURL     string    `json:"image"`
	ID           string    `json:"id"`
}

func (s *Service) PopulateProducts(ctx context.Context) error {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://62daf70dd1d97b9e0c49ca5d.mockapi."+
		"io/v1/products", nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.New("received error from remote server")
	}

	remoteProducts := []RemoteProduct{}
	err = json.NewDecoder(resp.Body).Decode(&remoteProducts)
	if err != nil {
		return err
	}

	for _, remoteProduct := range remoteProducts {
		product := app.Product{
			Name:               remoteProduct.ProductName,
			Description:        "",
			Image:              remoteProduct.ImageURL,
			AvailabilityStatus: "",
			AvailableItems:     0,
			CreatedAt:          remoteProduct.CreatedAt,
			Manufacturer:       remoteProduct.Manufacturer,
			Vehicle:            remoteProduct.Vehicle,
			ID:                 remoteProduct.ID,
		}
		err := s.repo.AddProduct(ctx, product)

		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
	}

	return nil
}

// ReviewOrderProducts requests from user to review the purchased products
func (s *Service) ReviewOrderProducts(ctx context.Context, conn *websocket.Conn, order *app.Order,
	orderProducts []app.OrderProduct) error {

	//Start discussion
	welcomeMessage := "Hey " + order.Customer.FirstName + "! "
	welcomeMessage += "Hope you have received your order you placed at " + order.PlacedDate.String() + " as expected! "
	welcomeMessage += "We would love some feedback for the products you have received! "
	if err := conn.WriteMessage(websocket.TextMessage, []byte(welcomeMessage)); err != nil {
		s.logger.With("success", false, "err", err)
		return err
	}

	for _, orderProduct := range orderProducts {
		askForProductReviewMessage := "Could you please share your experience with your purchase of " + orderProduct.
			Product.Name + "?"
		if err := conn.WriteMessage(websocket.TextMessage, []byte(askForProductReviewMessage)); err != nil {
			s.logger.With("success", false, "err", err)
			return err
		}

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			s.logger.With("success", false, "err", err)
			return err
		}
		analysisScore, err := s.sentimentAnalyzer.Process(ctx, string(p))
		if err != nil {
			s.logger.With("success", false, "err", err)
			return err
		}

		if err = s.repo.AddOrderProductReviewByOrderProductUUID(ctx, orderProduct.UUID,
			analysisScore.SentimentScore); err != nil {
			s.logger.With("success", false, "err", err)
			return err
		}

		generatedResponse, err := s.responseGenerator.Generate(ctx, analysisScore, orderProduct.Product.Name)
		if err != nil {
			s.logger.With("success", false, "err", err)
			return err
		}
		if err := conn.WriteMessage(messageType, []byte(generatedResponse.Response)); err != nil {
			s.logger.With("success", false, "err", err)
			return err
		}
	}
	err := s.repo.UpdateOrderStatusByOrderUUID(ctx, order.UUID, string(app.OrderStatusReviewed))
	if err != nil {
		s.logger.With("success", false, "err", err)
		return err
	}

	//End discussion
	thanksMessage := "Thank you for your time reviewing your products! "
	thanksMessage += "Hope to see you again " + order.Customer.LastName + "!"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(thanksMessage)); err != nil {
		s.logger.With("success", false, "err", err)
		return err
	}

	return nil
}
