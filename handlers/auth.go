package handlers

import (
	"encoding/json"
	"net/http"

	"go-tunnel/models"
	"go-tunnel/services"
	"go-tunnel/utils"
)

// AdminRegister handles admin registration
func AdminRegister(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, false, "Invalid JSON", nil)
		return
	}
	if req.Email == "" || req.Password == "" {
		utils.SendResponse(w, false, "Email va password kerak", nil)
		return
	}
	// TODO: Save admin credentials securely (hash password, store in DB or file)
	utils.SendResponse(w, true, "Admin ro'yxatdan o'tdi", map[string]string{"email": req.Email})
}

// AdminLogin handles admin login
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, false, "Invalid JSON", nil)
		return
	}
	if req.Email == "" || req.Password == "" {
		utils.SendResponse(w, false, "Email va password kerak", nil)
		return
	}
	// TODO: Check admin credentials (compare with stored hash)
	utils.SendResponse(w, true, "Admin login muvaffaqiyatli", map[string]interface{}{"email": req.Email})
}

// AdminGenerateToken generates a token for admin
func AdminGenerateToken(w http.ResponseWriter, r *http.Request) {
	// Generate token using service
	token := services.GenerateAdminToken()
	utils.SendResponse(w, true, "Token yaratildi", map[string]string{"token": token})
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.M
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.SendResponse(w, false, "Invalid JSON", nil)
        return
    }

    if req.Email == "" || req.Password == "" {
        utils.SendResponse(w, false, "Email va password kerak", nil)
        return
    }

    // Example data
    data := map[string]string{"email": req.Email}
    utils.SendResponse(w, true, "Login muvaffaqiyatli", data)
}

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, false, "Invalid JSON", nil)
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.SendResponse(w, false, "Email va password kerak", nil)
		return
	}

	data := map[string]interface{}{
		"email": req.Email,
		"token": "example_token_12345",
	}
	utils.SendResponse(w, true, "Login muvaffaqiyatli", data)
}

// Register handler
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, false, "Invalid JSON", nil)
		return
	}

	data := map[string]string{"email": req.Email}
	utils.SendResponse(w, true, "Ro'yxatdan o'tish muvaffaqiyatli", data)
}

// Logout handler
func Logout(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, true, "Logout muvaffaqiyatli", nil)
}

// RefreshToken handler
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"token": "new_token_12345"}
	utils.SendResponse(w, true, "Token yangilandi", data)
}

// HealthCheck handler
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok"}
	utils.SendResponse(w, true, "Server ishlayapti", data)
}

