package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	domain   = os.Getenv("AUTH0_DOMAIN")
	audience = os.Getenv("AUTH0_AUDIENCE")
)

// AuthMiddleware Middleware to validate JWT token.
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	jwksURL := fmt.Sprintf("%s.well-known/jwks.json", domain)

	ctx := context.Background()
	set, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		panic("failed to fetch JWKS: " + err.Error())
	}

	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, err := validateToken(tokenString, set)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !hasAllowedRole(claims, allowedRoles) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient privileges"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("missing or invalid Authorization header")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func validateToken(tokenString string, set jwk.Set) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		keyID, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing key ID (kid) in token header")
		}

		key, found := set.LookupKeyID(keyID)
		if !found {
			return nil, errors.New("unable to find matching key")
		}

		var rawKey interface{}
		if err := key.Raw(&rawKey); err != nil {
			return nil, errors.New("failed to parse key")
		}
		return rawKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func hasAllowedRole(claims jwt.MapClaims, allowedRoles []string) bool {
	roles, ok := claims["https://parkiu.com/roles"]
	if !ok {
		return false
	}

	switch v := roles.(type) {
	case string:
		return contains(allowedRoles, v)
	case []interface{}:
		for _, role := range v {
			if roleStr, ok := role.(string); ok && contains(allowedRoles, roleStr) {
				return true
			}
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
