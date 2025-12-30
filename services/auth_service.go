package services

import "go-tunnel/utils"

// GenerateAdminToken generates a secure token for admin
func GenerateAdminToken() string {
	return "admin_token_" + utils.GenerateRandomString(16)
}
