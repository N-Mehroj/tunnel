# EXAMPLES - Laravel-Style Routing va Middleware

## 1. Oddiy Route'ni Yaratish

```go
router := routes.NewRouter()

// GET request
router.GET("/users", handlers.GetUsers)

// POST request
router.POST("/users", handlers.CreateUser)

// PUT request
router.PUT("/users/{id}", handlers.UpdateUser)

// DELETE request
router.DELETE("/users/{id}", handlers.DeleteUser)

// Istalgan metodga
router.Any("/status", handlers.Status)
```

## 2. Middleware'ni Single Route'ga Qo'sh

```go
// Bitta middleware
router.GET("/users", handlers.GetUsers, middleware.CheckAuth)

// Bir nechta middleware
router.POST("/users", handlers.CreateUser, 
    middleware.CheckAuth,
    middleware.CheckContentType("application/json"),
    middleware.ErrorHandler)
```

## 3. Route Group (Prefix va Shared Middleware)

```go
// Barcha API route'lari /api prefix bilan va Auth bilan
router.Group("/api", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsers)
    group.POST("/users", handlers.CreateUser)
    group.GET("/posts", handlers.GetPosts)
}, middleware.CheckAuth)
```

## 4. Admin Routes

```go
// Admin-only route'lar /admin prefix bilan
router.Group("/admin", func(group *RouteGroup) {
    group.GET("/dashboard", handlers.AdminDashboard)
    group.GET("/logs", handlers.GetLogs)
    group.DELETE("/users/{id}", handlers.DeleteUser)
}, middleware.AdminOnly)
```

## 5. Public routes (Middleware'siz)

```go
router.POST("/auth/login", handlers.Login)
router.POST("/auth/register", handlers.Register)
router.GET("/health", handlers.HealthCheck)
```

## MIDDLEWARE EXAMPLES

### 1. CheckAuth Middleware

```go
// Authorization header'i zarur
router.GET("/profile", handlers.GetProfile, middleware.CheckAuth)

// Request:
// GET /profile
// Header: Authorization: Bearer token123
```

### 2. AdminOnly Middleware

```go
// Admin token'i zarur
router.POST("/users/{id}/delete", handlers.DeleteUser, middleware.AdminOnly)

// Request:
// POST /users/1/delete
// Header: X-Admin-Token: admin123
```

### 3. CORS Middleware

```go
// CORS header'larini qo'sh
router.Group("/api", func(group *RouteGroup) {
    group.GET("/data", handlers.GetData)
}, middleware.CORS)

// Javob beradi:
// Access-Control-Allow-Origin: *
// Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
```

### 4. Custom Middleware Yaratish

```go
// middleware/custom.go
package middleware

func CheckApiKey(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if apiKey != "valid-key-123" {
            http.Error(w, `{"error":"Invalid API key"}`, 
                http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

func CustomLogMiddleware(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request received: %s %s", r.Method, r.URL.Path)
        start := time.Now()
        
        next(w, r)
        
        log.Printf("Response time: %v", time.Since(start))
    }
}

// routes.go'da ishlatish:
router.GET("/secure-data", handlers.GetSecureData, 
    middleware.CheckApiKey)
```

## REAL WORLD EXAMPLES

### Example 1: E-Commerce API

```go
func SetupRoutes() *Router {
    router := NewRouter()

    // Public routes
    router.GET("/products", handlers.ListProducts)
    router.GET("/products/{id}", handlers.GetProduct)
    router.POST("/auth/login", handlers.Login)
    router.POST("/auth/register", handlers.Register)

    // User routes (Auth required)
    router.Group("/user", func(group *RouteGroup) {
        group.GET("/profile", handlers.GetProfile)
        group.PUT("/profile", handlers.UpdateProfile)
        group.GET("/orders", handlers.GetMyOrders)
        group.POST("/cart", handlers.AddToCart)
        group.POST("/checkout", handlers.Checkout)
    }, middleware.CheckAuth)

    // Admin routes (Admin only)
    router.Group("/admin", func(group *RouteGroup) {
        group.GET("/products", handlers.AdminListProducts)
        group.POST("/products", handlers.CreateProduct)
        group.PUT("/products/{id}", handlers.UpdateProduct)
        group.DELETE("/products/{id}", handlers.DeleteProduct)
        group.GET("/orders", handlers.AdminGetOrders)
        group.GET("/users", handlers.AdminListUsers)
    }, middleware.AdminOnly)

    return router
}
```

### Example 2: Blog API with Categories

```go
func SetupRoutes() *Router {
    router := NewRouter()

    // Public routes
    router.GET("/posts", handlers.ListPosts)
    router.GET("/posts/{id}", handlers.GetPost)
    router.GET("/categories", handlers.ListCategories)

    // Author routes (Auth required)
    router.Group("/author", func(group *RouteGroup) {
        group.POST("/posts", handlers.CreatePost)
        group.PUT("/posts/{id}", handlers.UpdatePost)
        group.DELETE("/posts/{id}", handlers.DeletePost)
        group.GET("/my-posts", handlers.GetMyPosts)
    }, middleware.CheckAuth)

    // Admin routes
    router.Group("/admin", func(group *RouteGroup) {
        group.POST("/categories", handlers.CreateCategory)
        group.PUT("/categories/{id}", handlers.UpdateCategory)
        group.DELETE("/categories/{id}", handlers.DeleteCategory)
        group.GET("/reports", handlers.GetReports)
    }, middleware.AdminOnly)

    return router
}
```

### Example 3: Multi-level Middleware

```go
// Her request'da run bo'ladigan middleware'lar
router.Group("/api", func(group *RouteGroup) {
    
    // Public sub-routes
    group.GET("/health", handlers.Health)
    
    // Protected sub-routes
    group.Group("/v1", func(subgroup *RouteGroup) {
        subgroup.GET("/data", handlers.GetData)
        subgroup.POST("/data", handlers.CreateData)
    }, middleware.CheckAuth)
    
}, middleware.CORS, middleware.LogRequest)
```

## Handler Implementation Examples

### Simple Handler

```go
func GetUsers(w http.ResponseWriter, r *http.Request) {
    users := []map[string]interface{}{
        {"id": 1, "name": "Ali"},
        {"id": 2, "name": "Zahra"},
    }
    utils.SendResponse(w, true, "Success", users)
}
```

### Handler with Request Body

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        utils.SendResponse(w, false, "Invalid JSON", nil)
        return
    }
    
    data := map[string]interface{}{
        "id":    1,
        "name":  user.Name,
        "email": user.Email,
    }
    utils.SendResponse(w, true, "User created", data)
}
```

### Handler with URL Parameter

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    // URL parameter'ni qabul qilish uchun chi?
    // Handler'da r.URL.Path'dan parse qilishingiz mumkin
    
    userId := "1" // Parse qiling URL'dan
    
    data := map[string]interface{}{
        "id":   userId,
        "name": "Ali",
    }
    utils.SendResponse(w, true, "User found", data)
}
```

## Testing with cURL

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"123456"}'

# Get users (with auth)
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer token123"

# Create user (with auth)
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer token123" \
  -H "Content-Type: application/json" \
  -d '{"name":"Nuri","email":"nuri@example.com"}'

# Admin dashboard
curl -X GET http://localhost:8080/admin/dashboard \
  -H "X-Admin-Token: admin123"

# Delete user (admin)
curl -X DELETE http://localhost:8080/admin/users/1 \
  -H "X-Admin-Token: admin123"
```

## Tips

1. **Middleware'lar ketma-ketda bajariladi** - eng oxirgi birinchi
2. **Group middleware'lari barcha ichidagi route'larga tadbiq etiladi**
3. **Handler type** - `func(http.ResponseWriter, *http.Request)` oddiy
4. **ConvertMiddleware** - middleware.Handler -> routes.Middleware uchun
5. **Path'ni tekshiring** - `/api/users` yoki `/api/users/{id}` larni to'g'ri yozing

---

**Keyingi qadamlar**: Database integration, JWT auth, validation qo'shin!
