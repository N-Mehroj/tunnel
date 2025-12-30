# Ngrok bilan TCP va HTTP API Testing Qo'llanmasi

## Qism 1: API Ro'yxati

Sizning Go server quyidagi API endpoints ni taqdim etadi:

### Ochiq API'lar (Authentication siz)
```
POST   /login              - Login endpoint
GET    /health             - Server health check
```

### Auth API'lar (Authentication kerak)
```
POST   /auth/register      - Foydalanuvchi ro'yxatdan o'tish
POST   /auth/login         - Login
POST   /auth/logout        - Logout (token kerak)
POST   /auth/refresh       - Token yangilash (token kerak)
```

### User API'lar (Auth middleware bilan)
```
GET    /api/users          - Barcha foydalanuvchilarni olish
POST   /api/users          - Yangi foydalanuvchi yaratish
GET    /api/users/{id}     - ID bo'yicha foydalanuvchi olish
PUT    /api/users/{id}     - Foydalanuvchini yangilash
DELETE /api/users/{id}     - Foydalanuvchini o'chirish
```

### Admin API'lar (Admin token kerak)
```
GET    /admin/dashboard    - Admin dashboard
GET    /admin/users        - Barcha foydalanuvchilar (admin)
DELETE /admin/users/{id}   - Foydalanuvchini o'chirish (admin)
```

---

## Qism 2: Ngrok o'rnatish va Ishlatish

### 2.1 Ngrok o'rnatish

**MacOS (Homebrew):**
```bash
brew install ngrok
```

**Linux:**
```bash
# O'rnatish
wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip
unzip ngrok-stable-linux-amd64.zip
sudo mv ngrok /usr/local/bin/

# Tekshirish
ngrok version
```

**Windows (PowerShell):**
```powershell
choco install ngrok
# yoki
scoop install ngrok
```

### 2.2 Ngrok autentifikatsiya

```bash
ngrok authtoken YOUR_TOKEN_HERE
```

[Ngrok account yaratish uchun](https://dashboard.ngrok.com/signup)

### 2.3 Server Ishga Tushirish

```bash
# Go server ishga tushirish (port 8080 da)
go run main.go

# Yoki
APP_URL=0.0.0.0:8080 go run main.go
```

### 2.4 Ngrok Tunnel Ochish

```bash
# HTTP tunnel (API testing uchun)
ngrok http 8080

# Chiqish:
# ngrok                                                         (Ctrl+C to quit)
# Forwarding                 http://abc123def456.ngrok.io -> http://localhost:8080
```

**MUHIM:** URL `http://abc123def456.ngrok.io` dan client uchun ishlatish kerak!

---

## Qism 3: HTTP API Client (Python)

### 3.1 Requirements o'rnatish

```bash
pip install requests
```

### 3.2 Test Script yaratish

`test_api_client.py` fayl yarating:

```python
import requests
import json

# Ngrok URL (o'z URL ingizni o'rnatish kerak!)
BASE_URL = "http://abc123def456.ngrok.io"

class APIClient:
    def __init__(self, base_url):
        self.base_url = base_url
        self.token = None
        
    def register(self, email, password):
        """Ro'yxatdan o'tish"""
        url = f"{self.base_url}/auth/register"
        data = {"email": email, "password": password}
        response = requests.post(url, json=data)
        print(f"✓ Register: {response.status_code}")
        print(f"  Response: {response.json()}\n")
        return response.json()
    
    def login(self, email, password):
        """Login"""
        url = f"{self.base_url}/auth/login"
        data = {"email": email, "password": password}
        response = requests.post(url, json=data)
        result = response.json()
        if result.get('success'):
            self.token = result.get('data', {}).get('token')
        print(f"✓ Login: {response.status_code}")
        print(f"  Response: {result}\n")
        return result
    
    def get_users(self):
        """Barcha foydalanuvchilarni olish"""
        url = f"{self.base_url}/api/users"
        headers = self._get_headers()
        response = requests.get(url, headers=headers)
        print(f"✓ Get Users: {response.status_code}")
        print(f"  Response: {response.json()}\n")
        return response.json()
    
    def create_user(self, email, password):
        """Yangi foydalanuvchi yaratish"""
        url = f"{self.base_url}/api/users"
        headers = self._get_headers()
        data = {"email": email, "password": password}
        response = requests.post(url, json=data, headers=headers)
        print(f"✓ Create User: {response.status_code}")
        print(f"  Response: {response.json()}\n")
        return response.json()
    
    def get_user(self, user_id):
        """ID bo'yicha foydalanuvchi olish"""
        url = f"{self.base_url}/api/users/{user_id}"
        headers = self._get_headers()
        response = requests.get(url, headers=headers)
        print(f"✓ Get User {user_id}: {response.status_code}")
        print(f"  Response: {response.json()}\n")
        return response.json()
    
    def health_check(self):
        """Server holatini tekshirish"""
        url = f"{self.base_url}/health"
        response = requests.get(url)
        print(f"✓ Health Check: {response.status_code}")
        print(f"  Response: {response.json()}\n")
        return response.json()
    
    def _get_headers(self):
        """Authorization header yaratish"""
        headers = {"Content-Type": "application/json"}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        return headers

# Test Script
if __name__ == "__main__":
    client = APIClient(BASE_URL)
    
    print("=" * 60)
    print("GO HTTP API CLIENT TEST")
    print("=" * 60 + "\n")
    
    # 1. Health check
    print("1. Server holatini tekshirish:")
    client.health_check()
    
    # 2. Ro'yxatdan o'tish
    print("2. Foydalanuvchi ro'yxatdan o'tish:")
    client.register("testuser@example.com", "password123")
    
    # 3. Login
    print("3. Login qilish:")
    client.login("testuser@example.com", "password123")
    
    # 4. Barcha foydalanuvchilarni olish
    print("4. Barcha foydalanuvchilarni olish:")
    client.get_users()
    
    # 5. Yangi foydalanuvchi yaratish
    print("5. Yangi foydalanuvchi yaratish:")
    client.create_user("newuser@example.com", "pass456")
    
    # 6. Foydalanuvchini olish
    print("6. ID bo'yicha foydalanuvchi olish:")
    client.get_user(1)
    
    print("=" * 60)
    print("TEST YAKUNLANDI")
    print("=" * 60)
```

### 3.3 Test qilish

```bash
python test_api_client.py
```

---

## Qism 4: TCP Client Connection (Advanced)

Agar siz native TCP connection qilmoqchi bo'lsangiz (HTTP emas):

### 4.1 Go TCP Server yaratish

`tcp_server.go` fayl yarating:

```go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// TCPMessage - TCP orqali yuborish uchun struktura
type TCPMessage struct {
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Data   map[string]string `json:"data,omitempty"`
	Token  string            `json:"token,omitempty"`
}

// TCPResponse - TCP javob struktura
type TCPResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	
	reader := bufio.NewReader(conn)
	
	for {
		// JSON message o'qish
		data, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		
		var msg TCPMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			response := TCPResponse{
				Status:  400,
				Success: false,
				Message: "Invalid JSON format",
			}
			respondTCP(conn, response)
			continue
		}
		
		// Request qayta ishlab chiqish
		response := processTCPRequest(msg)
		respondTCP(conn, response)
	}
}

func processTCPRequest(msg TCPMessage) TCPResponse {
	switch msg.Method {
	case "LOGIN":
		return TCPResponse{
			Status:  200,
			Success: true,
			Message: "Login successful",
			Data: map[string]string{
				"email": msg.Data["email"],
				"token": "tcp_token_12345",
			},
		}
	case "REGISTER":
		return TCPResponse{
			Status:  201,
			Success: true,
			Message: "User registered successfully",
		}
	case "GET_USERS":
		return TCPResponse{
			Status:  200,
			Success: true,
			Message: "Users retrieved",
			Data: []map[string]string{
				{"id": "1", "email": "user1@example.com"},
				{"id": "2", "email": "user2@example.com"},
			},
		}
	default:
		return TCPResponse{
			Status:  400,
			Success: false,
			Message: "Unknown method: " + msg.Method,
		}
	}
}

func respondTCP(conn net.Conn, response TCPResponse) {
	data, _ := json.Marshal(response)
	conn.Write(append(data, '\n'))
}

func startTCPServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		fmt.Println("TCP Server xatosi:", err)
		return
	}
	defer listener.Close()
	
	fmt.Println("TCP Server 9000 portda ishga tushdi...")
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection xatosi:", err)
			continue
		}
		go handleTCPConnection(conn)
	}
}

func main() {
	startTCPServer()
}
```

### 4.2 Python TCP Client

`tcp_client.py` fayl yarating:

```python
import socket
import json

class TCPClient:
    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.socket = None
        
    def connect(self):
        """Serverga ulanish"""
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.socket.connect((self.host, self.port))
            print(f"✓ {self.host}:{self.port} ga ulandi\n")
        except Exception as e:
            print(f"✗ Connection xatosi: {e}")
            
    def send_request(self, method, data=None, token=None):
        """Request yuborish"""
        message = {
            "method": method,
            "path": "",
            "data": data or {},
            "token": token
        }
        
        try:
            self.socket.send((json.dumps(message) + '\n').encode())
            response_data = self.socket.recv(4096).decode()
            response = json.loads(response_data)
            
            print(f"→ Yuborildi: {method}")
            print(f"← Javob: {json.dumps(response, indent=2)}\n")
            return response
        except Exception as e:
            print(f"✗ Request xatosi: {e}")
            
    def close(self):
        """Ulanishni yopish"""
        if self.socket:
            self.socket.close()
            print("Ulanish yopildi")

# Test
if __name__ == "__main__":
    print("=" * 60)
    print("TCP CLIENT TEST")
    print("=" * 60 + "\n")
    
    client = TCPClient("localhost", 9000)
    client.connect()
    
    # Test requests
    client.send_request("REGISTER", {"email": "test@example.com", "password": "123"})
    client.send_request("LOGIN", {"email": "test@example.com", "password": "123"})
    client.send_request("GET_USERS")
    
    client.close()
```

### 4.3 TCP Test qilish

**Terminal 1 - Server:**
```bash
go run tcp_server.go
# Output: TCP Server 9000 portda ishga tushdi...
```

**Terminal 2 - Client:**
```bash
python tcp_client.py
```

---

## Qism 5: Ngrok bilan Ishlarni Qo'llash

### Kombinlangan Setup:

```bash
# Terminal 1: Go Server
go run main.go

# Terminal 2: Ngrok
ngrok http 8080

# Terminal 3: Client Test
# test_api_client.py dagi BASE_URL ni ngrok URL ga o'rgating
# python test_api_client.py
```

### cURL bilan Test

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

# Get Users (token kerak)
curl -X GET http://abc123def456.ngrok.io/api/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Qism 6: Debugging Tips

### Logs ko'rish
```bash
# Ngrok logs
tail -f ~/.ngrok2/ngrok.log

# Go server debug
APP_URL=0.0.0.0:8080 go run main.go
```

### Common Xatolar

| Xato | Yechim |
|------|--------|
| `Connection refused` | Server ishga tushmaganini tekshir |
| `Invalid token` | Token mavjudligini tekshir |
| `CORS error` | CORS middleware tekshir |
| `Ngrok tunnel expired` | Yangi tunnel ochish kerak |

---

## Xulosa

- **HTTP API**: Ngrok orqali public URL
- **TCP Connection**: Port 9000 orqali native connection
- **Python Clients**: HTTP va TCP test scriptlari ready
- **All APIs**: Login, Register, CRUD operations

Fikrlar yoki savollarniz bo'lsa, help ga murojaat qiling! 🚀
