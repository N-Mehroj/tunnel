package handlers

import (
	"encoding/json"
	"net/http"

	"go-tunnel/utils"
)

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []map[string]interface{}{
		{"id": 1, "name": "Ali", "email": "ali@example.com"},
		{"id": 2, "name": "Zahra", "email": "zahra@example.com"},
	}
	utils.SendResponse(w, true, "Foydalanuvchilar ro'yxati", users)
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendResponse(w, false, "Invalid JSON", nil)
		return
	}

	data := map[string]interface{}{
		"id":   3,
		"name": user["name"],
		"email": user["email"],
	}
	utils.SendResponse(w, true, "Foydalanuvchi yaratildi", data)
}

// GetUser returns a single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"id": 1,
		"name": "Ali",
		"email": "ali@example.com",
	}
	utils.SendResponse(w, true, "Foydalanuvchi topildi", data)
}

// UpdateUser updates a user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendResponse(w, false, "Invalid JSON", nil)
		return
	}

	data := map[string]interface{}{
		"id": 1,
		"name": user["name"],
		"email": user["email"],
	}
	utils.SendResponse(w, true, "Foydalanuvchi yangilandi", data)
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, true, "Foydalanuvchi o'chirildi", nil)
}

// GetAllUsers returns all users (admin)
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := []map[string]interface{}{
		{"id": 1, "name": "Ali", "email": "ali@example.com", "role": "user"},
		{"id": 2, "name": "Zahra", "email": "zahra@example.com", "role": "user"},
		{"id": 3, "name": "Admin", "email": "admin@example.com", "role": "admin"},
	}
	utils.SendResponse(w, true, "Barcha foydalanuvchilar", users)
}

// DeleteUserAdmin deletes a user (admin only)
func DeleteUserAdmin(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, true, "Foydalanuvchi admin tomonidan o'chirildi", nil)
}

// AdminDashboard shows admin dashboard
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"total_users": 3,
		"total_requests": 1000,
		"status": "online",
	}
	utils.SendResponse(w, true, "Admin boshqaruvi paneli", data)
}
