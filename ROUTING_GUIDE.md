# Laravel-Style Routing Guide

Bu qo'llanma Go uygulamasi uchun Laravel-uslubidagi marshrutlash tizimini tushuntirib beradi.

## Fayllar Tuzilishi

```
routes/
├── router.go          # Router asosiy strukturasi
└── routes.go          # Barcha marshrutlar

middleware/
└── middleware.go      # Middleware funksiyalari

handlers/
├── auth.go           # Authentication handlerlari
└── users.go          # User handlerlari

main.go               # Ilovani boshlash
```

## Asosiy Tushunchalar

### 1. Router - Laravel Router'iga o'xshash

```go
router := routes.NewRouter()

// GET, POST, PUT, DELETE, PATCH, Any
router.GET("/users", handlers.GetUsers)
router.POST("/users", handlers.CreateUser)
```

### 2. Middleware - Requestlarni filtrlash

Middleware'lar requestni interceptlar va ishlaydi:

```go
router.GET("/users", handlers.GetUsers, middleware.CheckAuth)
```

### 3. Route Groups - Marshrutlarni guruhlash

Laravel'dagi `Route::group()` kabi:

```go
router.Group("/api", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsers)
    group.POST("/users", handlers.CreateUser)
}, middleware.CheckAuth)
```

## Misol: Hozirgi Routinglar

### Public Routes (Ochiq)
```
POST /login           - Tizimga kirish
GET  /health          - Server holati
```

### Auth Routes (CORS bilan)
```
POST /auth/register   - Ro'yxatdan o'tish
POST /auth/login      - Tizimga kirish
POST /auth/logout     - Tizimdan chiqish (Auth zarur)
POST /auth/refresh    - Tokenni yangilash (Auth zarur)
```

### API Routes (Auth zarur)
```
GET    /api/users         - Barcha foydalanuvchilar
POST   /api/users         - Yangi foydalanuvchi qo'sh
GET    /api/users/{id}    - Bitta foydalanuvchi
PUT    /api/users/{id}    - Foydalanuvchini yangilash
DELETE /api/users/{id}    - Foydalanuvchini o'chirish
```

### Admin Routes (Admin zarur)
```
GET    /admin/dashboard       - Admin boshqaruvi paneli
GET    /admin/users           - Barcha foydalanuvchilar (admin ko'rinish)
DELETE /admin/users/{id}      - Foydalanuvchini o'chirish (admin)
```

## Middleware'lar

### Mavjud Middleware'lar

| Middleware | Maqsad | Misol |
|-----------|--------|-------|
| `CheckAuth` | Authorization headerini tekshirish | `middleware.CheckAuth` |
| `AdminOnly` | Faqat admin uchun | `middleware.AdminOnly` |
| `CORS` | CORS headerlarini qo'sh | `middleware.CORS` |
| `CheckContentType` | Content-Type tekshirish | `middleware.CheckContentType("application/json")` |
| `CheckMethod` | HTTP metodni tekshirish | `middleware.CheckMethod("GET", "POST")` |
| `ErrorHandler` | Xatolarni boshqarish | `middleware.ErrorHandler` |
| `LogRequest` | Requestlarni logga yozish | `middleware.LogRequest` |
| `RateLimit` | Chastota cheklash | `middleware.RateLimit(100)` |

## Middleware'ni Foydalanish

### 1. Bitta Route'ga Middleware

```go
router.GET("/users", handlers.GetUsers, middleware.CheckAuth)
```

### 2. Bir nechta Middleware'ni Ketma-ket

```go
router.GET("/users", handlers.GetUsers, 
    middleware.CheckAuth,
    middleware.CORS,
    middleware.LogRequest)
```

### 3. Group'ga Middleware

```go
router.Group("/api", func(group *RouteGroup) {
    group.GET("/users", handlers.GetUsers)
}, middleware.CheckAuth, middleware.CORS)
```

## O'z Middleware'ini Yaratish

```go
// middleware/custom.go
package middleware

func CheckApiKey(next Handler) Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        if apiKey == "" {
            http.Error(w, `{"error":"API key required"}`, 
                http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}
```

Keyin routes'da:

```go
router.POST("/api-call", handlers.ApiCall, middleware.CheckApiKey)
```

## Handler Yaratish

Handler - oddiy funksiya:

```go
package handlers

import (
    "net/http"
    "go-tunnel/utils"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
    // Request bilan ishlash
    data := map[string]interface{}{
        "message": "Salom",
    }
    utils.SendResponse(w, true, "Muvaffaqiyatli", data)
}
```

## Testing

### cURL orqali Test

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"123456"}'

# Users (Auth kerak)
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer token123"

# Admin (Admin kerak)
curl -X GET http://localhost:8080/admin/dashboard \
  -H "X-Admin-Token: admin123"
```

## Xatolar va Ruxsat Berish

### 401 Unauthorized
```
Header: Authorization zarur emas
Status: 401
Body: {"error":"Unauthorized"}
```

### 403 Forbidden
```
Header: X-Admin-Token zarur emas
Status: 403
Body: {"error":"Admin access required"}
```

### 405 Method Not Allowed
```
Noto'g'ri HTTP method
Status: 405
Body: {"error":"Method not allowed"}
```

## Keyingi Qadamlar

1. **Database Integration** - Handlers'da database query'larini qo'shin
2. **JWT Auth** - Token bazali autentifikatsiya qo'shin
3. **Validation** - Request parametrlarini tekshiring
4. **Error Handling** - O'z error handlerlari yarating
5. **Logging** - Batafsil logga yozish qo'shin

## Tips

- Middleware'lar **ketma-ketda** bajariladi (eng oxirgi birinchi)
- Group middleware'lari barcha ichidagi route'larga qo'llaniladi
- `Handler` - oddiy `func(http.ResponseWriter, *http.Request)` turida
- Route path'ini tekshirish uchun `/path` dan boshlang

---

**Savollar bo'lsa, handler yoki middleware'ni yangilang va test qiling!**
