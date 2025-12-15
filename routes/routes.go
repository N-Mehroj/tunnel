package routes

import (
	"net/http"

	"go-tunnel/handlers"
	"go-tunnel/middleware"
)

// ConvertMiddleware converts middleware.Handler to routes.Middleware
func ConvertMiddleware(m func(middleware.Handler) middleware.Handler) Middleware {
	return func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			// Convert h to middleware.Handler
			mh := middleware.Handler(h)
			// Apply middleware
			result := m(mh)
			// Call result with ResponseWriter and Request
			result(w, r)
		}
	}
}

// SetupRoutes configures all application routes with middleware (Laravel-style)
func SetupRoutes() *Router {
	router := NewRouter()

	// Public routes (without middleware)
	router.POST("/login", handlers.LoginHandler)
	router.GET("/health", handlers.HealthCheck)

	// API routes with auth middleware
	router.Group("/api", func(group *RouteGroup) {
		// User routes with auth middleware
		group.GET("/users", handlers.GetUsers, ConvertMiddleware(middleware.CheckAuth))
		group.POST("/users", handlers.CreateUser, ConvertMiddleware(middleware.CheckAuth))
		group.GET("/users/{id}", handlers.GetUser, ConvertMiddleware(middleware.CheckAuth))
		group.PUT("/users/{id}", handlers.UpdateUser, ConvertMiddleware(middleware.CheckAuth))
		group.DELETE("/users/{id}", handlers.DeleteUser, ConvertMiddleware(middleware.CheckAuth))
	}, ConvertMiddleware(middleware.CheckAuth), ConvertMiddleware(middleware.CORS))

	// Admin routes with admin middleware
	router.Group("/admin", func(group *RouteGroup) {
		group.GET("/dashboard", handlers.AdminDashboard, ConvertMiddleware(middleware.AdminOnly))
		group.GET("/users", handlers.GetAllUsers, ConvertMiddleware(middleware.AdminOnly))
		group.DELETE("/users/{id}", handlers.DeleteUserAdmin, ConvertMiddleware(middleware.AdminOnly))
	}, ConvertMiddleware(middleware.AdminOnly))

	// Auth routes
	router.Group("/auth", func(group *RouteGroup) {
		group.POST("/register", handlers.Register)
		group.POST("/login", handlers.Login)
		group.POST("/logout", handlers.Logout, ConvertMiddleware(middleware.CheckAuth))
		group.POST("/refresh", handlers.RefreshToken, ConvertMiddleware(middleware.CheckAuth))
	}, ConvertMiddleware(middleware.CORS))

	return router
}

