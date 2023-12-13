package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"reviewbot/app"
	"time"
)

// CustomerResponse represents an order response object entity.
type CustomerResponse struct {
	UUID             string    `json:"uuid"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	PhoneNumber      string    `json:"phone_number"`
	RegistrationDate time.Time `json:"registration_date"`
}

// OrderResponse represents an order response object entity.
type OrderResponse struct {
	UUID       string           `json:"uuid"`
	Customer   CustomerResponse `json:"customer"`
	Status     app.OrderStatus  `json:"status"`
	PlacedDate time.Time        `json:"placed_date"`
}

// OrderStatusRequest represents an order status update request object entity.
type OrderStatusRequest struct {
	Status app.OrderStatus `json:"status"`
}

// OrderProductResponse represents an order products object entity.
type OrderProductResponse struct {
	UUID        string          `json:"uuid"`
	OrderUUID   string          `json:"order_uuid"`
	ProductUUID string          `json:"product_uuid"`
	Items       int             `json:"items"`
	Product     ProductResponse `json:"product"`
}

// ProductResponse represents a product object entity.
type ProductResponse struct {
	UUID               string `json:"uuid"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Image              string `json:"items"`
	AvailabilityStatus string `json:"availability_status"`
	AvailableItems     int    `json:"available_items"`
}

func (srv *Server) getOrderByUUID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), handlerDefaultTimeout)
	defer cancel()
	log := srv.App.Logger.With(LogFieldKeyRequestID, GetReqID(ctx))

	orderUUID := mux.Vars(r)["order_uuid"]
	order, err := srv.UserService.OrderByUUID(ctx, orderUUID)
	if err != nil {
		log.Error(err.Error())
		log.With("success", false, "err", err)
		if errors.Is(err, app.ErrNoRecords) {
			NotFoundError(w, err)
			return
		}
		ServerError(w, err)
		return
	}

	srv.App.Logger.With("success", true)
	Ok(w, transformOrderToResponse(order), http.StatusOK)
}

func (srv *Server) updateOrderStatusByUUID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), handlerDefaultTimeout)
	defer cancel()
	log := srv.App.Logger.With(LogFieldKeyRequestID, GetReqID(ctx))

	orderUUID := mux.Vars(r)["order_uuid"]
	orderStatusRequest := OrderStatusRequest{}
	err := json.NewDecoder(r.Body).Decode(&orderStatusRequest)
	if err != nil {
		log.Error(err.Error())
		log.With("success", false, "err", err)
		BadRequestError(w, err)
		return
	}
	if orderStatusRequest.Status != app.OrderStatusPreparing && orderStatusRequest.Status != app.OrderStatusPlaced &&
		orderStatusRequest.Status != app.OrderStatusSending && orderStatusRequest.Status != app.OrderStatusCompleted {
		err = errors.New("invalid status provided")
		log.With("success", false, "err", err)
		BadRequestError(w, err)
		return
	}

	err = srv.UserService.UpdateOrderStatusByUUID(ctx, orderUUID, orderStatusRequest.Status)
	if err != nil {
		log.Error(err.Error())
		log.With("success", false, "err", err)
		if errors.Is(err, app.ErrNoRecords) {
			NotFoundError(w, err)
			return
		}
		ServerError(w, err)
		return
	}
	srv.App.Logger.With("success", true)
	Ok(w, nil, http.StatusNoContent)
}

func (srv *Server) getOrderProductsByOrderUUID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), handlerDefaultTimeout)
	defer cancel()
	log := srv.App.Logger.With(LogFieldKeyRequestID, GetReqID(ctx))

	orderUUID := mux.Vars(r)["order_uuid"]
	orderProducts, err := srv.UserService.OrderProductsByOrderUUID(ctx, orderUUID)
	if err != nil {
		log.Error(err.Error())
		log.With("success", false, "err", err)
		if errors.Is(err, app.ErrNoRecords) {
			NotFoundError(w, err)
			return
		}
		ServerError(w, err)
		return
	}

	srv.App.Logger.With("success", true)
	Ok(w, transformOrderProductsToResponse(orderProducts), http.StatusOK)
}

func transformOrderToResponse(order *app.Order) OrderResponse {
	return OrderResponse{
		UUID: order.UUID,
		Customer: CustomerResponse{
			UUID:             order.Customer.UUID,
			FirstName:        order.Customer.FirstName,
			LastName:         order.Customer.LastName,
			Email:            order.Customer.Email,
			PhoneNumber:      order.Customer.PhoneNumber,
			RegistrationDate: order.Customer.RegistrationDate,
		},
		Status:     order.Status,
		PlacedDate: order.PlacedDate,
	}
}

func transformOrderProductsToResponse(orderProducts []app.OrderProduct) []OrderProductResponse {
	orderProductResponse := []OrderProductResponse{}
	for _, orderProduct := range orderProducts {
		orderProductResponse = append(orderProductResponse, OrderProductResponse{
			UUID:        orderProduct.UUID,
			OrderUUID:   orderProduct.OrderUUID,
			ProductUUID: orderProduct.ProductUUID,
			Items:       orderProduct.Items,
			Product: ProductResponse{
				UUID:               orderProduct.Product.UUID,
				Name:               orderProduct.Product.Name,
				Description:        orderProduct.Product.Description,
				Image:              orderProduct.Product.Image,
				AvailabilityStatus: orderProduct.Product.AvailabilityStatus,
				AvailableItems:     orderProduct.Product.AvailableItems,
			},
		})
	}
	return orderProductResponse
}
