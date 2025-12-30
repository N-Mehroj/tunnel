# 🚀 Ngrok va TCP Testing - Tez Qo'llanma

## 📋 Ishni Ketma-Ketligini

### **Variant 1: HTTP API Testing (Ngrok orqali)**

#### Terminal 1 - Go Server ishga tushirish
```bash
cd /home/nmehroj/Desktop/go/http
go run main.go

# Chiqish: Server running at http://0.0.0.0:8080
```

#### Terminal 2 - Ngrok tunnel ochish
```bash
ngrok http 8080

# Chiqish (URL ni copy qil):
# Forwarding    http://abc123def456.ngrok.io -> http://localhost:8080
```

#### Terminal 3 - Python client test qilish (Lokal)
```bash
cd /home/nmehroj/Desktop/go/http
python test_api_client.py
```

### **Variant 2: TCP Native Connection**

#### Terminal 1 - TCP Server ishga tushirish
```bash
cd /home/nmehroj/Desktop/go/http
go run tcp_server.go

# Chiqish: TCP SERVER ISHGA TUSHDI, Port: 9000
```

#### Terminal 2 - TCP Client test qilish
```bash
cd /home/nmehroj/Desktop/go/http
python tcp_client.py
```

---

## 🌐 Ngrok orqali API Testing

### Step 1: BASE_URL o'rgat
`test_api_client.py` dagi 12-qatorni o'rgat:
```python
# Eski:
BASE_URL = "http://localhost:8080"

# Yangi (Ngrok URL):
BASE_URL = "http://abc123def456.ngrok.io"  # O'z URL ingizni qo'ying!
```

### Step 2: Test qil
```bash
python test_api_client.py
```

### Step 3: cURL bilan ham test qil
```bash
# Health Check
curl -X GET http://abc123def456.ngrok.io/health

# Register
curl -X POST http://abc123def456.ngrok.io/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass123"}'

# Login
curl -X POST http://abc123def456.ngrok.io/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass123"}'
```

---

## 📡 TCP Connection Tafsiloti

### TCP Server Methods
```
LOGIN          - Login  
REGISTER       - Ro'yxatdan o'tish
GET_USERS      - Barcha users
GET_USER       - Bitta user (id kerak)
CREATE_USER    - Yangi user
DELETE_USER    - User o'chirish
```

### TCP Request Format
```json
{
  "method": "LOGIN",
  "data": {
    "email": "test@example.com",
    "password": "password123"
  },
  "token": ""
}
```

### TCP Response Format
```json
{
  "status": 200,
  "success": true,
  "message": "Login muvaffaqiyatli",
  "data": {
    "email": "test@example.com",
    "token": "tcp_token_test@example.com"
  }
}
```

---

## 📂 Fayllar

| Fayl | Maqsad |
|------|--------|
| `test_api_client.py` | HTTP API client test |
| `tcp_server.go` | TCP server (Port 9000) |
| `tcp_client.py` | TCP client test |
| `NGROK_TCP_SETUP.md` | To'liq qo'llanma |

---

## 🔧 Debugging

### Problem: Connection refused
```bash
# Serverning ishga tushganini tekshir
netstat -tuln | grep 8080    # HTTP
netstat -tuln | grep 9000    # TCP
```

### Problem: Ngrok xatosi
```bash
# Token o'rnatish
ngrok authtoken YOUR_TOKEN_HERE

# Token tekshirish
ngrok version
```

### Problem: Python xatosi
```bash
# Requests o'rnatish
pip install requests

# Python version
python --version
```

---

## 💡 Foydalı Linklar

- [Ngrok Dashboard](https://dashboard.ngrok.com)
- [Go Documentation](https://golang.org/doc)
- [REST API Best Practices](https://restfulapi.net)
- [TCP Socket Programming](https://golang.org/pkg/net)

---

## ⚠️ Muhim Eslatmalar

1. **Ngrok URL** har 2 soatda o'zgaradi (free plan)
2. **Token** middleware dan keyin yuboriladi
3. **TCP** HTTP dan boshqachroq (no HTTP headers)
4. **CORS** middleware orqali enable qilingan
5. **Auth** `/api` va `/admin` routes uchun kerak

---

**Tayyor! Ngrok bilan testing qilishni boshlashingiz mumkin! 🚀**
