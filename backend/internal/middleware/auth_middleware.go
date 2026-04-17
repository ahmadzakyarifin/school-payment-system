package middleware

import (
	"strings"

	"github.com/ahmadzakyarifin/school-payment-system/pkg/response"
	"github.com/ahmadzakyarifin/school-payment-system/pkg/utils/token"
	"github.com/gin-gonic/gin"
)

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Token diperlukan")
			c.Abort()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Format token salah (harus Bearer <token>)")
			c.Abort()
			return
		}

		claims, err := token.Validate(parts[1], secret)
		if err != nil {
			response.Unauthorized(c, "Token tidak valid atau kedaluwarsa")
			c.Abort()
			return
		}

		c.Set("user_id", claims["sub"])
		c.Set("email", claims["email"]) 
		c.Set("role", claims["role"])

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Unauthorized(c, "Informasi user tidak ditemukan")
			c.Abort()
			return
		}
		roleAllowed := false
		for _, r := range roles {
			if r == userRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			response.Forbidden(c, "Anda tidak memiliki izin untuk akses fitur ini")
			c.Abort()
			return
		}

		c.Next()
	}
}
