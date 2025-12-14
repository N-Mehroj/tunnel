package utils

import (
	"encoding/json"
	"net/http"

	"go-tunnel/models"
)

func SendResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    resp := models.Response{
        Success: success,
        Message: message,
        Data:    data,
    }
    json.NewEncoder(w).Encode(resp)
}
