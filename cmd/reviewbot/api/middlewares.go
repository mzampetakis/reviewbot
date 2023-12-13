package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if w.Header().Get("Content-Type") == "" {
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				}
				w.WriteHeader(http.StatusInternalServerError)
				f := "PANIC: "
				stack := fmt.Sprint(f, err, ", stack trace:", '\n')
				app.Logger.Error(stack)
				fmt.Println(string(debug.Stack()))
				stack = fmt.Sprint(stack, string(debug.Stack()))
				w.Write([]byte(stack))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

// ContextKey is
type ContextKey string

const (

	// ContextKeyReqID is the context key for RequestID
	ContextKeyReqID      ContextKey = "requestID"
	LogFieldKeyRequestID string     = string(ContextKeyReqID)
	// HTTPHeaderNameRequestID has the name of the header for request ID
	HTTPHeaderNameRequestID = "X-Request-ID"
)

// GetReqID will get reqID from a http request and return it as a string
func GetReqID(ctx context.Context) string {
	reqID := ctx.Value(ContextKeyReqID)
	if ret, ok := reqID.(string); ok {
		return ret
	}
	return ""
}

// AttachReqID will attach a new request ID to a http request
func AttachReqID(ctx context.Context) context.Context {
	reqID := uuid.New().String()
	return context.WithValue(ctx, ContextKeyReqID, reqID)
}

func (app *Application) httpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := AttachReqID(r.Context())
		r = r.WithContext(ctx)
		app.Logger.Info(fmt.Sprintf("Request %s-start: %s %s", GetReqID(ctx), r.Method, r.URL.Path))
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		responseBody := ""
		if lrw.statusCode >= http.StatusBadRequest {
			responseBody = " Response body: " + string(lrw.responseData) + "."
		}
		app.Logger.Info(fmt.Sprintf("Completed request %s-end: %s %s %v (%v) in %v.%s", GetReqID(ctx), r.Method,
			r.URL.Path, http.StatusText(lrw.statusCode), lrw.statusCode, time.Since(start), responseBody))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseData []byte
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, []byte{}}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseData = b
	return size, err
}
