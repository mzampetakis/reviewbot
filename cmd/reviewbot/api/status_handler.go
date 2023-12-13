package api

import (
	"net/http"
)

func (srv *Server) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "OK",
	}
	err := JSON(w, http.StatusOK, data)
	if err != nil {
		srv.App.serverError(w, r, err)
	}
}
