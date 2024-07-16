package handler

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        if token == "" || !strings.HasPrefix(token, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
            c.Abort()
            return
        }

        actualToken := strings.TrimPrefix(token, "Bearer ")

        if !validateToken(actualToken) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func validateToken(token string) bool {
    return token == "ABC123456"
}
