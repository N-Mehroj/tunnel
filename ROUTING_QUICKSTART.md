# Routing va Middleware - Quick Start

Go Laravel-uslubida routing va middleware bilan tez boshlash qo'llanmasi.

## 📋 Fayl Tuzilishi

```
routes/
  ├── router.go      # Router class (Laravel Router'iga o'xshash)
  └── routes.go      # Barcha routes

middleware/
  └── middleware.go  # Middleware funksiyalari

handlers/
  ├── auth.go        # Auth handlers
  └── users.go       # User handlers

main.go            # Entry point
```

## 🚀 Asosiy Foydalanish

### 1. Router Yaratish va Routes Qo'sh

```go
// routes/routes.go
func SetupRoutes() *Router {
    router := NewRouter()
    
    // GET route
    router.GET("/home", handlers.Home)
    
    // POST route
    router.POST("/login", handlers.Login)
    
    // Middleware bilan
    router.GET("/profile", handlers.Profile, middleware.CheckAuth)
    
    return router
}
```

### 2. Main'da Router'ni Ishlatish

```go
// main.go
func main() {
    router := routes.SetupRoutes()
    
    http.ListenAndServe(":8080", router)
}
```

### 3. Handler Yaratish

```go
// handlers/home.go
func Home(w http.ResponseWriter, r *http.Request) {
    utils.SendResponse(w, true, "Salom!", nil)
}
```

## 🎯 Laravel-Style Route Methods

```go
router.GET(path, handler)          // GET
router.POST(path, handler)         // POST
router.PUT(path, handler)          // PUT
router.DELETE(path, handler)       // DELETE
router.PATCH(path, handler)        // PATCH
router.Any(path, handler)          // Barcha metodlar
```

## 🔐 Middleware

### Mavjud Middleware'lar

| Middleware | Maqsad |
|-----------|--------|
| `CheckAuth` | Authorization header tekshiring |
| `AdminOnly` | Faqat admin uchun |
| `CORS` | CORS enable qiling |
| `ErrorHandler` | Xatolarni tutun |
| `LogRequest` | Request'larni log qiling |

### Middleware'ni Ishlatish

```go
// Bitta middleware
router.GET("/users", handlers.GetUsers, middleware.CheckAuth)

// Bir nechta middleware
router.POST("/users", handlers.CreateUser, 
    middleware.CheckAuth, 
    middleware.ErrorHandler)
```

## 📦 Route Groups (Grouping)

```go
router.Group("/api", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsers)
    group.POST("/users", handlers.CreateUser)
}, middleware.CheckAuth)

// Bu /api/users ga POST request'da middleware.CheckAuth ishlatiladi
```

## 💡 Amaliy Misol

```go
func SetupRoutes() *Router {
    router := NewRouter()

    // Public
    router.GET("/health", handlers.HealthCheck)
    router.POST("/login", handlers.Login)

    // User routes (Auth required)
    router.Group("/user", func(group *RouteGroup) {
        group.GET("/profile", handlers.GetProfile)
        group.PUT("/profile", handlers.UpdateProfile)
    }, middleware.CheckAuth)

    // Admin routes (Admin only)
    router.Group("/admin", func(group *RouteGroup) {
        group.GET("/users", handlers.AdminListUsers)
        group.DELETE("/users/{id}", handlers.DeleteUser)
    }, middleware.AdminOnly)

    return router
}
```

## 🧪 Testing

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

# Protected route (with auth)
curl http://localhost:8080/user/profile \
  -H "Authorization: Bearer token123"

# Admin route (with admin token)
curl http://localhost:8080/admin/users \
  -H "X-Admin-Token: admin123"
```

## 📝 O'z Middleware'ini Yaratish

```go
// middleware/middleware.go
func MyMiddleware(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        // Before handler
        log.Println("Before handler")
        
        // Call handler
        next(w, r)
        
        // After handler
        log.Println("After handler")
    }
}

// routes.go'da
router.GET("/page", handlers.MyHandler, middleware.MyMiddleware)
```

## 🔍 Debugging

Routes'ni tekshirish:

```go
router := routes.SetupRoutes()
for _, route := range router.ListRoutes() {
    fmt.Printf("%s %s\n", route.Method, route.Path)
}
```

## ✅ Checklist

- [ ] Router'ni main.go'da ishlatish
- [ ] Public routes qo'sh
- [ ] Auth routes qo'sh
- [ ] Handler'larni yaratish
- [ ] Middleware qo'sh
- [ ] cURL bilan test qilish
- [ ] Database integration qo'sh
- [ ] Validation qo'sh

---

**Keyingi**: ROUTING_GUIDE.md - Batafsil qo'llanma
