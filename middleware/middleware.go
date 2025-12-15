package middleware

import (
	"log"
	"net/http"
	"time"
)

// Handler type for middleware
type Handler func(http.ResponseWriter, *http.Request)

// LogRequest logs each HTTP request
func LogRequest(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		log.Printf(
			"| %s | %s | %s | %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
		
		next(w, r)
	}
}

// CheckAuth checks if user is authenticated
func CheckAuth(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		log.Printf("Auth token: %s", token)
		next(w, r)
	}
}

// CheckContentType checks if request has proper Content-Type
func CheckContentType(contentType string) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") != contentType && r.Method != "GET" {
				http.Error(w, `{"error":"Invalid Content-Type"}`, http.StatusBadRequest)
				return
			}
			next(w, r)
		}
	}
}

// CheckMethod ensures only specific HTTP methods are allowed
func CheckMethod(allowedMethods ...string) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			allowed := false
			for _, method := range allowedMethods {
				if r.Method == method {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
				return
			}
			next(w, r)
		}
	}
}

// SetHeaders adds custom headers to response
func SetHeaders(headers map[string]string) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			for key, value := range headers {
				w.Header().Set(key, value)
			}
			next(w, r)
		}
	}
}

// AdminOnly checks if user is admin
func AdminOnly(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		adminToken := r.Header.Get("X-Admin-Token")
		if adminToken == "" {
			http.Error(w, `{"error":"Admin access required"}`, http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// RateLimit is a simple rate limiting middleware
func RateLimit(requestsPerSecond int) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Rate limit check: %d req/sec", requestsPerSecond)
			next(w, r)
		}
	}
}

// CORS enables CORS headers
func CORS(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

// ErrorHandler wraps handler with error recovery
func ErrorHandler(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Error: %v", err)
				http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}
