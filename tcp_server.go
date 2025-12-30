package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// TCPMessage - TCP orqali yuborish uchun struktura
type TCPMessage struct {
	Method string            `json:"method"`
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

// HandleTCPConnection - Client connection qayta ishlash
func handleTCPConnection(conn net.Conn, clientID int) {
	defer conn.Close()
	
	fmt.Printf("[Client %d] Ulandi: %s\n", clientID, conn.RemoteAddr())
	
	reader := bufio.NewReader(conn)
	
	for {
		// JSON message o'qish
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[Client %d] Disconnected\n", clientID)
			break
		}
		
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		var msg TCPMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Printf("[Client %d] Invalid JSON: %v\n", clientID, err)
			response := TCPResponse{
				Status:  400,
				Success: false,
				Message: "Invalid JSON format",
			}
			respondTCP(conn, response)
			continue
		}
		
		fmt.Printf("[Client %d] Method: %s\n", clientID, msg.Method)
		
		// Request qayta ishlab chiqish
		response := processTCPRequest(msg)
		respondTCP(conn, response)
	}
}

// ProcessTCPRequest - Request qayta ishlash
func processTCPRequest(msg TCPMessage) TCPResponse {
	switch strings.ToUpper(msg.Method) {
	case "LOGIN":
		email := msg.Data["email"]
		password := msg.Data["password"]
		
		if email == "" || password == "" {
			return TCPResponse{
				Status:  400,
				Success: false,
				Message: "Email va password kerak",
			}
		}
		
		return TCPResponse{
			Status:  200,
			Success: true,
			Message: "Login muvaffaqiyatli",
			Data: map[string]string{
				"email": email,
				"token": "tcp_token_" + email,
			},
		}
	
	case "REGISTER":
		email := msg.Data["email"]
		password := msg.Data["password"]
		
		if email == "" || password == "" {
			return TCPResponse{
				Status:  400,
				Success: false,
				Message: "Email va password kerak",
			}
		}
		
		return TCPResponse{
			Status:  201,
			Success: true,
			Message: "Foydalanuvchi muvaffaqiyatli ro'yxatdan o'ttdi",
			Data: map[string]string{
				"email": email,
			},
		}
	
	case "GET_USERS":
		return TCPResponse{
			Status:  200,
			Success: true,
			Message: "Foydalanuvchilar olindi",
			Data: []map[string]string{
				{"id": "1", "email": "user1@example.com"},
				{"id": "2", "email": "user2@example.com"},
				{"id": "3", "email": "user3@example.com"},
			},
		}
	
	case "GET_USER":
		userID := msg.Data["id"]
		if userID == "" {
			return TCPResponse{
				Status:  400,
				Success: false,
				Message: "User ID kerak",
			}
		}
		
		return TCPResponse{
			Status:  200,
			Success: true,
			Message: "Foydalanuvchi olindi",
			Data: map[string]string{
				"id":    userID,
				"email": "user" + userID + "@example.com",
			},
		}
	
	case "CREATE_USER":
		email := msg.Data["email"]
		password := msg.Data["password"]
		
		if email == "" || password == "" {
			return TCPResponse{
				Status:  400,
				Success: false,
				Message: "Email va password kerak",
			}
		}
		
		return TCPResponse{
			Status:  201,
			Success: true,
			Message: "Foydalanuvchi yaratildi",
			Data: map[string]string{
				"id":    "4",
				"email": email,
			},
		}
	
	case "DELETE_USER":
		userID := msg.Data["id"]
		if userID == "" {
			return TCPResponse{
				Status:  400,
				Success: false,
				Message: "User ID kerak",
			}
		}
		
		return TCPResponse{
			Status:  200,
			Success: true,
			Message: "Foydalanuvchi o'chirildi",
			Data: map[string]string{
				"deleted_id": userID,
			},
		}
	
	default:
		return TCPResponse{
			Status:  400,
			Success: false,
			Message: "Noto'g'ri method: " + msg.Method,
		}
	}
}

// RespondTCP - TCP javob yuborish
func respondTCP(conn net.Conn, response TCPResponse) {
	data, _ := json.Marshal(response)
	conn.Write(append(data, '\n'))
}

// StartTCPServer - TCP Server ishga tushirish
func startTCPServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Printf("TCP Server xatosi: %v\n", err)
		return
	}
	defer listener.Close()
	
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println(" " + strings.Repeat(" ", 18) + "TCP SERVER ISHGA TUSHDI")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("\n📡 Port: 9000")
	fmt.Println("🔌 Address: 127.0.0.1:9000")
	fmt.Println("\n⏳ Client lari kutilmoqda...\n")
	fmt.Println(strings.Repeat("=", 70) + "\n")
	
	clientID := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Connection xatosi: %v\n", err)
			continue
		}
		clientID++
		go handleTCPConnection(conn, clientID)
	}
}

func main() {
	startTCPServer()
}
