package routes

import (
	"net/http"
)

// Handler is the type for route handlers
type Handler func(http.ResponseWriter, *http.Request)

// Middleware is the type for middleware functions - accepts Handler type
type Middleware func(Handler) Handler

// RouteGroup represents a group of routes with shared prefix and middleware
type RouteGroup struct {
	prefix     string
	middleware []Middleware
	routes     []Route
}

// Route represents a single route
type Route struct {
	Method     string
	Path       string
	Handler    Handler
	Middleware []Middleware
}

// Router is the main router structure (like Laravel Router)
type Router struct {
	routes  []Route
	groups  []RouteGroup
	baseMux *http.ServeMux
}

// NewRouter creates a new router instance
func NewRouter() *Router {
	return &Router{
		routes:  []Route{},
		groups:  []RouteGroup{},
		baseMux: http.NewServeMux(),
	}
}

// GET registers a GET route
func (r *Router) GET(path string, handler Handler, middleware ...Middleware) *Router {
	r.addRoute("GET", path, handler, middleware...)
	return r
}

// POST registers a POST route
func (r *Router) POST(path string, handler Handler, middleware ...Middleware) *Router {
	r.addRoute("POST", path, handler, middleware...)
	return r
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler Handler, middleware ...Middleware) *Router {
	r.addRoute("PUT", path, handler, middleware...)
	return r
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler Handler, middleware ...Middleware) *Router {
	r.addRoute("DELETE", path, handler, middleware...)
	return r
}

// PATCH registers a PATCH route
func (r *Router) PATCH(path string, handler Handler, middleware ...Middleware) *Router {
	r.addRoute("PATCH", path, handler, middleware...)
	return r
}

// Any registers a route for any HTTP method
func (r *Router) Any(path string, handler Handler, middleware ...Middleware) *Router {
	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH"} {
		r.addRoute(method, path, handler, middleware...)
	}
	return r
}

// Group creates a route group with shared prefix and middleware
func (r *Router) Group(prefix string, groupFunc func(group *RouteGroup), middleware ...Middleware) *Router {
	group := &RouteGroup{
		prefix:     prefix,
		middleware: middleware,
		routes:     []Route{},
	}
	groupFunc(group)
	
	// Register all routes in the group
	for _, route := range group.routes {
		fullPath := group.prefix + route.Path
		// Combine group middleware with route middleware
		allMiddleware := append(group.middleware, route.Middleware...)
		r.addRoute(route.Method, fullPath, route.Handler, allMiddleware...)
	}
	
	return r
}

// Group methods for RouteGroup (to allow nested route definitions)
func (rg *RouteGroup) GET(path string, handler Handler, middleware ...Middleware) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:     "GET",
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	})
	return rg
}

func (rg *RouteGroup) POST(path string, handler Handler, middleware ...Middleware) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:     "POST",
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	})
	return rg
}

func (rg *RouteGroup) PUT(path string, handler Handler, middleware ...Middleware) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:     "PUT",
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	})
	return rg
}

func (rg *RouteGroup) DELETE(path string, handler Handler, middleware ...Middleware) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:     "DELETE",
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	})
	return rg
}

func (rg *RouteGroup) PATCH(path string, handler Handler, middleware ...Middleware) *RouteGroup {
	rg.routes = append(rg.routes, Route{
		Method:     "PATCH",
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	})
	return rg
}

// addRoute adds a route with middleware to the router
func (r *Router) addRoute(method string, path string, handler Handler, middleware ...Middleware) {
	// Apply middleware to handler
	finalHandler := handler
	for i := len(middleware) - 1; i >= 0; i-- {
		finalHandler = middleware[i](finalHandler)
	}

	// Store route info
	r.routes = append(r.routes, Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	})

	// Register the path if not already registered
	alreadyRegistered := false
	for _, route := range r.routes[:len(r.routes)-1] {
		if route.Path == path {
			alreadyRegistered = true
			break
		}
	}

	if !alreadyRegistered {
		// Create a dispatcher that checks method
		r.baseMux.HandleFunc(path, r.dispatcher)
	}
}

// dispatcher is a custom dispatcher that matches method and calls appropriate handler
func (r *Router) dispatcher(w http.ResponseWriter, req *http.Request) {
	// Find matching route by path and method
	for _, route := range r.routes {
		if route.Path == req.URL.Path && (route.Method == req.Method || route.Method == "ANY") {
			// Apply middleware
			handler := route.Handler
			for i := len(route.Middleware) - 1; i >= 0; i-- {
				handler = route.Middleware[i](handler)
			}
			handler(w, req)
			return
		}
	}
	// No matching route found
	http.NotFound(w, req)
}

// ServeHTTP makes Router implement http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.baseMux.ServeHTTP(w, req)
}

// ListRoutes returns all registered routes (for debugging)
func (r *Router) ListRoutes() []Route {
	return r.routes
}
