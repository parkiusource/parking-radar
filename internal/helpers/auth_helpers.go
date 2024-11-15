package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ExtractAdminIDAndRole extracts the admin ID and role from the JWT token.
func ExtractAdminIDAndRole(c *gin.Context) (string, bool) {
	userClaims, ok := c.Get("user")
	if !ok {
		return "", false
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	adminID, _ := claims["sub"].(string)

	roles, ok := claims["https://parkiu.com/roles"].([]interface{})
	if !ok {
		return adminID, false
	}

	for _, role := range roles {
		if roleStr, ok := role.(string); ok && roleStr == "admin_global" {
			return adminID, true
		}
	}
	return adminID, false
}
