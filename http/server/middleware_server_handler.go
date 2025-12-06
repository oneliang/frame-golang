package server

import (
	"github.com/oneliang/frame-golang/http/middleware"
	"net/http"
)

// MiddlewareServerHandler .
type MiddlewareServerHandler struct {
	middleware    []middleware.Middleware
	serverHandler *ServerHandler
}

// NewMiddlewareServerHandler .
func NewMiddlewareServerHandler(serverHandler *ServerHandler) *MiddlewareServerHandler {
	return &MiddlewareServerHandler{
		middleware: make([]middleware.Middleware, 0),
	}
}

// ServeHTTP .
func (this *MiddlewareServerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var handler http.Handler = http.HandlerFunc(this.serverHandler.ServeHTTP)
	for i := len(this.middleware) - 1; i >= 0; i-- {
		handler = this.middleware[i](handler)
	}
	handler.ServeHTTP(writer, request)
}

// AddMiddleware .
func (this *MiddlewareServerHandler) AddMiddleware(middleware middleware.Middleware) {
	this.middleware = append(this.middleware, middleware)
}
