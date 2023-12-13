package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (srv *Server) routes() http.Handler {
	serverMux := mux.NewRouter()

	serverMux.NotFoundHandler = http.HandlerFunc(srv.App.notFound)
	serverMux.MethodNotAllowedHandler = http.HandlerFunc(srv.App.methodNotAllowed)
	serverMux.Use(srv.App.httpLogger)
	serverMux.Use(srv.App.recoverPanic)
	serverMux.Use(srv.App.enableCORS)

	apiMux := serverMux.PathPrefix("/api").Subrouter()
	apiMux.HandleFunc("/status", srv.status).Methods("GET")

	ordersMux := apiMux.PathPrefix("/orders").Subrouter()
	ordersMux.HandleFunc("/{order_uuid}", srv.getOrderByUUID).Methods("GET")
	ordersMux.HandleFunc("/{order_uuid}", srv.updateOrderStatusByUUID).Methods("PATCH")
	ordersMux.HandleFunc("/{order_uuid}/products", srv.getOrderProductsByOrderUUID).Methods("GET")

	return serverMux
}
