package middleware

import "net/http"

// Middleware .
type Middleware func(http.Handler) http.Handler

// Handle .
func (this Middleware) Handle(next http.Handler) http.Handler {
	return this(next)
}
