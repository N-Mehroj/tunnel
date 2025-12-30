package utils

import (
	"encoding/json"
	"go-tunnel/models"
	"math/rand"
	"net/http"
	"time"
)

// GenerateRandomString returns a random string of given length
func GenerateRandomString(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}



func SendResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    resp := models.Response{
        Success: success,
        Message: message,
        Data:    data,
    }
    json.NewEncoder(w).Encode(resp)
}
