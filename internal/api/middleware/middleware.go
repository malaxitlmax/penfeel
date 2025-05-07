package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pb "github.com/malaxitlmax/penfeel/api/proto"
	"golang.org/x/net/context"
)

// AuthMiddleware middleware для авторизации через JWT токен
func AuthMiddleware(authClient pb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Проверяем формат Bearer token
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format, should be 'Bearer {token}'"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(splitToken[1])
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token cannot be empty"})
			c.Abort()
			return
		}

		// Проверяем токен через auth service
		res, err := authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
			Token: token,
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !res.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": res.Error})
			c.Abort()
			return
		}

		// Сохраняем информацию о пользователе в контексте
		c.Set("user_id", res.User.Id)
		c.Set("username", res.User.Username)
		c.Set("email", res.User.Email)

		c.Next()
	}
}
