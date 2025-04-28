package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/malaxitlmax/penfeel/api/proto"
	"golang.org/x/net/context"
)

// Handler структура обработчика авторизации
type Handler struct {
	authClient pb.AuthServiceClient
}

// NewHandler создает новый обработчик авторизации
func NewHandler(authClient pb.AuthServiceClient) *Handler {
	return &Handler{
		authClient: authClient,
	}
}

// RegisterRequest структура запроса на регистрацию
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest структура запроса на вход
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenRequest структура запроса на валидацию токена
type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// Register обрабатывает запрос на регистрацию
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправляем запрос к auth-сервису через gRPC
	res, err := h.authClient.Register(context.Background(), &pb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user_id": res.UserId,
	})
}

// Login обрабатывает запрос на вход
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправляем запрос к auth-сервису через gRPC
	res, err := h.authClient.Login(context.Background(), &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !res.Success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         res.Token,
		"refresh_token": res.RefreshToken,
		"user": gin.H{
			"id":       res.User.Id,
			"username": res.User.Username,
			"email":    res.User.Email,
		},
	})
}

// ValidateToken обрабатывает запрос на валидацию токена
func (h *Handler) ValidateToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Отправляем запрос к auth-сервису через gRPC
	res, err := h.authClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: req.Token,
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if !res.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user": gin.H{
			"id":       res.User.Id,
			"username": res.User.Username,
			"email":    res.User.Email,
		},
	})
}
