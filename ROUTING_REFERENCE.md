# Laravel-Style Routing - Visual Reference Card

## 🎯 Quick Reference

### Basic Routes

```go
router.GET("/path", handler)
router.POST("/path", handler)
router.PUT("/path", handler)
router.DELETE("/path", handler)
router.PATCH("/path", handler)
router.Any("/path", handler)
```

### With Middleware

```go
// Single middleware
router.GET("/path", handler, middleware.CheckAuth)

// Multiple middleware
router.POST("/path", handler, 
    middleware.CheckAuth,
    middleware.CheckContentType("application/json"),
    middleware.ErrorHandler)

// Group middleware
router.Group("/api", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsers)
}, middleware.CheckAuth)
```

---

## 📡 Request/Response Flow

```
┌─────────────────┐
│   Request       │
│  POST /users    │
│  with auth      │
└────────┬────────┘
         │
         ▼
    ┌─────────┐
    │ Router  │ ← Route match
    │ found   │
    └────┬────┘
         │
         ▼
    ┌──────────────────────┐
    │  Middleware Chain    │
    │                      │
    │  ConvertMiddleware   │
    │    ↓                 │
    │  CheckAuth           │
    │    ↓                 │
    │  ErrorHandler        │
    │    ↓                 │
    │  Handler             │
    └────┬─────────────────┘
         │
         ▼
    ┌───────────────────┐
    │  Response         │
    │  JSON data        │
    └───────────────────┘
```

---

## 🔒 Middleware Types

### Access Control

| Middleware | Purpose | Header | Example |
|-----------|---------|--------|---------|
| CheckAuth | Require auth | Authorization | Bearer token123 |
| AdminOnly | Require admin | X-Admin-Token | admin123 |

### Headers

| Middleware | Purpose | Adds |
|-----------|---------|------|
| CORS | Cross-origin | Access-Control-* headers |
| SetHeaders | Custom headers | User-defined headers |

### Validation

| Middleware | Purpose | Checks |
|-----------|---------|--------|
| CheckContentType | Content type | application/json |
| CheckMethod | HTTP method | GET, POST, etc. |

### Utilities

| Middleware | Purpose |
|-----------|---------|
| LogRequest | Log requests |
| ErrorHandler | Catch errors |
| RateLimit | Limit requests |

---

## 🛣️ Route Organization

### Pattern 1: Flat Routes

```go
router.GET("/", handlers.Home)
router.GET("/about", handlers.About)
router.POST("/contact", handlers.Contact)
```

### Pattern 2: Grouped Routes

```go
router.Group("/api", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsers)
    group.POST("/users", handlers.CreateUser)
}, middleware.CheckAuth)

// Results in:
// GET  /api/users
// POST /api/users
```

### Pattern 3: Mixed Routes

```go
// Public
router.GET("/home", handlers.Home)

// Protected
router.Group("/api", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsers)
}, middleware.CheckAuth)

// Admin only
router.Group("/admin", func(group *RouteGroup) {
    group.DELETE("/users/{id}", handlers.Delete)
}, middleware.AdminOnly)
```

---

## 📝 Handler Template

```go
package handlers

import (
    "encoding/json"
    "net/http"
    "go-tunnel/utils"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var req struct {
        Name string `json:"name"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.SendResponse(w, false, "Invalid JSON", nil)
        return
    }
    
    // Process
    data := map[string]interface{}{
        "name": req.Name,
    }
    
    // Response
    utils.SendResponse(w, true, "Success", data)
}
```

---

## 🔧 Custom Middleware Template

```go
// middleware/middleware.go

func MyMiddleware(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        // Before handler
        log.Println("Before")
        
        // Call next handler
        next(w, r)
        
        // After handler
        log.Println("After")
    }
}
```

Usage:
```go
router.GET("/path", handler, middleware.MyMiddleware)
```

---

## 🧪 Testing with cURL

### Public Route
```bash
curl http://localhost:8080/health
```

### Protected Route
```bash
curl -H "Authorization: Bearer token123" \
  http://localhost:8080/api/users
```

### POST Request
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'
```

### Admin Route
```bash
curl -H "X-Admin-Token: admin123" \
  http://localhost:8080/admin/dashboard
```

---

## 📊 Route Decision Tree

```
Is route public?
├─ YES → No middleware
└─ NO
   │
   ├─ Requires auth?
   │  ├─ YES → Add CheckAuth
   │  └─ NO
   │
   └─ Requires admin?
      ├─ YES → Add AdminOnly
      └─ NO → Custom middleware
```

---

## 🎯 Common Patterns

### Pattern: API with Versioning

```go
router.Group("/api/v1", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsersV1)
}, middleware.CheckAuth)

router.Group("/api/v2", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsersV2)
}, middleware.CheckAuth)
```

### Pattern: Nested Groups

```go
router.Group("/api", func(group *RouteGroup) {
    group.Group("/users", func(subgroup *RouteGroup) {
        subgroup.GET("", handlers.ListUsers)
        subgroup.POST("", handlers.CreateUser)
    }, middleware.CheckAuth)
})
```

### Pattern: Public + Protected

```go
router.POST("/auth/login", handlers.Login)

router.Group("/user", func(group *RouteGroup) {
    group.GET("/profile", handlers.GetProfile)
}, middleware.CheckAuth)
```

---

## 🔐 Security Checklist

- [ ] Auth routes protected with CheckAuth
- [ ] Admin routes protected with AdminOnly
- [ ] CORS configured correctly
- [ ] Input validation added
- [ ] Error messages don't leak info
- [ ] Sensitive data not logged
- [ ] Rate limiting considered
- [ ] HTTPS in production

---

## 📈 Performance Tips

1. **Minimize middleware** - Only use needed middleware
2. **Fast middleware** - Avoid heavy operations in middleware
3. **Cache responses** - Add caching middleware for GET requests
4. **Connection pooling** - Reuse DB connections
5. **Compress responses** - Add gzip middleware

---

## 🚨 Common Errors

| Error | Cause | Fix |
|-------|-------|-----|
| `Unauthorized` | Missing Authorization header | Add auth token |
| `Admin access required` | Missing X-Admin-Token | Add admin token |
| `Invalid JSON` | Bad request body | Check JSON format |
| `route conflict` | Duplicate routes | Check routes/routes.go |

---

## 📚 File Locations

| File | Purpose |
|------|---------|
| `routes/router.go` | Router implementation |
| `routes/routes.go` | Route configuration |
| `middleware/middleware.go` | Middleware functions |
| `handlers/` | Handler functions |
| `main.go` | Entry point |

---

## 🎓 Learning Resources

Inside your project:
1. **README_ROUTING.md** - Documentation index
2. **ROUTING_QUICKSTART.md** - Quick start
3. **ROUTING_GUIDE.md** - Detailed guide
4. **ROUTING_EXAMPLES.md** - Real examples

---

## ⚡ Quick Commands

```bash
# Build
go build -o app main.go

# Run
./app

# Test
curl http://localhost:8080/health

# View logs
tail -f /tmp/app.log
```

---

**Print this page for quick reference! 📄**
