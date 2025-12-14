package handlers

import (
	"encoding/json"
	"net/http"

	"go-tunnel/models"
	"go-tunnel/utils"
)

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
