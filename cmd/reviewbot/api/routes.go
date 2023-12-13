package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (srv *Server) routes() http.Handler {
	serverMux := mux.NewRouter()

	apiMux := serverMux.PathPrefix("/api").Subrouter()
	apiMux.NotFoundHandler = http.HandlerFunc(srv.App.notFound)
	apiMux.MethodNotAllowedHandler = http.HandlerFunc(srv.App.methodNotAllowed)
	apiMux.Use(srv.App.httpLogger)
	apiMux.Use(srv.App.recoverPanic)
	apiMux.Use(srv.App.enableCORS)
	apiMux.HandleFunc("/status", srv.status).Methods("GET")

	ordersMux := apiMux.PathPrefix("/orders").Subrouter()
	ordersMux.HandleFunc("/{order_uuid}", srv.getOrderByUUID).Methods("GET")
	ordersMux.HandleFunc("/{order_uuid}", srv.updateOrderStatusByUUID).Methods("PATCH")
	ordersMux.HandleFunc("/{order_uuid}/products", srv.getOrderProductsByOrderUUID).Methods("GET")

	wsMux := serverMux.PathPrefix("/ws").Subrouter()
	ordersWSMux := wsMux.PathPrefix("/orders").Subrouter()
	ordersWSMux.Handle("/{order_uuid}", srv)

	return serverMux
}
