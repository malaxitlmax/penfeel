package document

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pb "github.com/malaxitlmax/penfeel/api/proto"
	"golang.org/x/net/context"
)

// Handler структура обработчика документов
type Handler struct {
	documentClient pb.DocumentServiceClient
}

// NewHandler создает новый обработчик документов
func NewHandler(documentClient pb.DocumentServiceClient) *Handler {
	return &Handler{
		documentClient: documentClient,
	}
}

// GetDocumentsRequest структура запроса на получение документов
type GetDocumentsRequest struct {
	UserID string `json:"user_id" form:"user_id"`
}

// GetDocuments обрабатывает запрос на получение списка документов
func (h *Handler) GetDocuments(c *gin.Context) {
	var req GetDocumentsRequest

	// Получаем ID пользователя из запроса или из токена
	userID := c.GetString("user_id")
	if userID == "" {
		// Если ID не получен из токена, пробуем получить из query-параметров
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}
		userID = req.UserID
	}

	// Проверяем, что ID пользователя - валидный UUID
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	// Отправляем запрос к document-сервису через gRPC
	res, err := h.documentClient.GetDocuments(context.Background(), &pb.GetDocumentsRequest{
		UserId: userID,
	})

	if err != nil {
		// Проверяем ошибку на наличие "connection refused"
		if strings.Contains(err.Error(), "connection refused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "document service is unavailable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"documents": res.Documents,
	})
}

// GetDocument обрабатывает запрос на получение документа
func (h *Handler) GetDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document id is required"})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Отправляем запрос к document-сервису через gRPC
	res, err := h.documentClient.GetDocument(context.Background(), &pb.GetDocumentRequest{
		Id:     documentID,
		UserId: userID,
	})

	if err != nil {
		// Проверяем ошибку на наличие "connection refused"
		if strings.Contains(err.Error(), "connection refused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "document service is unavailable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"document": res.Document,
	})
}

// CreateDocumentRequest структура запроса на создание документа
type CreateDocumentRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// CreateDocument обрабатывает запрос на создание документа
func (h *Handler) CreateDocument(c *gin.Context) {
	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Отправляем запрос к document-сервису через gRPC
	res, err := h.documentClient.CreateDocument(context.Background(), &pb.CreateDocumentRequest{
		Title:   req.Title,
		Content: req.Content,
		UserId:  userID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "document service is unavailable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"document": res.Document,
	})
}

// UpdateDocumentRequest структура запроса на обновление документа
type UpdateDocumentRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// UpdateDocument обрабатывает запрос на обновление документа
func (h *Handler) UpdateDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document id is required"})
		return
	}

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Отправляем запрос к document-сервису через gRPC
	res, err := h.documentClient.UpdateDocument(context.Background(), &pb.UpdateDocumentRequest{
		Id:      documentID,
		Title:   req.Title,
		Content: req.Content,
		UserId:  userID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "document service is unavailable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"document": res.Document,
	})
}

// DeleteDocument обрабатывает запрос на удаление документа
func (h *Handler) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document id is required"})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Отправляем запрос к document-сервису через gRPC
	res, err := h.documentClient.DeleteDocument(context.Background(), &pb.DeleteDocumentRequest{
		Id:     documentID,
		UserId: userID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "document service is unavailable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// RegisterRoutes регистрирует маршруты для документов
func (h *Handler) RegisterRoutes(r *gin.Engine, middleware ...gin.HandlerFunc) {
	documents := r.Group("/api/documents")
	documents.Use(middleware...)

	documents.GET("", h.GetDocuments)
	documents.GET("/:id", h.GetDocument)
	documents.POST("", h.CreateDocument)
	documents.PUT("/:id", h.UpdateDocument)
	documents.DELETE("/:id", h.DeleteDocument)
}
