package middleware

import (
	"context"
	"net/http"
)

// contextServiceName is the type used for context keys
type contextServiceName string

// contextServiceNameKey is the key used in context for the service name
const contextServiceNameKey contextServiceName = "cacophony_service"

// Service middleware stores the given service name in the request context, it can be retriev later using GetService
func Service(service string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, contextServiceNameKey, service)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// GetService returns the service name of a given context, if none found it will return empty
func GetService(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	service, _ := ctx.Value(contextServiceNameKey).(string)
	return service
}
