package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"golang.org/x/exp/slog"
)

type malformedRequestError struct {
	status int
	msg    string
	err    error
}

func (mr *malformedRequestError) Error() string {
	return mr.msg
}

// ApiError defines the standard error the API returns.
type ApiError struct {
	Message string  `json:"message"`
	Errors  []error `json:"errors,omitempty"`
}

// ErrorWrapper wraps an ApiError, so consumers can check for it easier.
type ErrorWrapper struct {
	Error ApiError `json:"error"`
}

func (app *Application) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)
	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.Logger.Error(message, requestAttrs, "trace", trace)
}

func (app *Application) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string,
	headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]
	err := JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportServerError(r, err)
	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *Application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}
