package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	pb "github.com/malaxitlmax/penfeel/api/proto"
	"github.com/malaxitlmax/penfeel/internal/api-gateway/service"
	"golang.org/x/net/context"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Или реализуйте более строгую проверку происхождения
	},
}

// DocumentHandler структура обработчика документов
type DocumentHandler struct {
	documentClient pb.DocumentServiceClient
	wsService      *service.WebSocketService
}

// NewDocumentHandler создает новый обработчик документов
func NewDocumentHandler(documentClient pb.DocumentServiceClient) *DocumentHandler {
	return &DocumentHandler{
		documentClient: documentClient,
		wsService:      service.NewWebSocketService(documentClient),
	}
}

// GetDocumentsRequest структура запроса на получение документов
type GetDocumentsRequest struct {
	UserID string `json:"user_id" form:"user_id"`
}

// GetDocuments обрабатывает запрос на получение списка документов
func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	var req GetDocumentsRequest

	// Получаем ID пользователя из запроса или из токена
	userID := c.GetString("user_id")
	if userID == "" {
		// Если ID не получен из токена, пробуем получить из query-параметров
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: user_id is required"})
			return
		}
		userID = req.UserID
	}

	// Проверяем, что ID пользователя - валидный UUID
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format", "details": "User ID must be a valid UUID"})
		return
	}

	// Отправляем запрос к document-сервису через gRPC
	res, err := h.documentClient.GetDocuments(context.Background(), &pb.GetDocumentsRequest{
		UserId: userID,
	})

	if err != nil {
		// Проверяем ошибку на наличие "connection refused"
		if strings.Contains(err.Error(), "connection refused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Document service is unavailable - please try again later",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch documents",
			"details": err.Error(),
		})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Document service rejected the request",
			"details": res.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"documents": res.Documents,
	})
}

// GetDocument обрабатывает запрос на получение документа
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing document ID", "details": "Document ID is required in the path"})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required", "details": "Valid authentication token is required"})
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
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Document service is unavailable - please try again later",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Not Found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":      "Document not found",
				"details":    "The requested document doesn't exist or you don't have permission to access it",
				"debug_info": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch document",
			"details": err.Error(),
		})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Document service rejected the request",
			"details": res.Error,
		})
		return
	}

	// Только после успешного получения документа апгрейдим соединение до WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket", "details": err.Error()})
		return
	}

	// Передаем управление соединением в сервис WebSocket
	h.wsService.HandleWebSocketConnection(documentID, userID, conn, res.Document)
}

// CreateDocumentRequest структура запроса на создание документа
type CreateDocumentRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// CreateDocument обрабатывает запрос на создание документа
func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required: user ID missing from token"})
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
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Document service is unavailable - please try again later",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Not Found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":      "Failed to create document: resource not found",
				"details":    "This could be due to missing user permissions or service configuration issues",
				"debug_info": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create document",
			"details": err.Error(),
		})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Document service rejected the request",
			"details": res.Error,
		})
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
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing document ID", "details": "Document ID is required in the path"})
		return
	}

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required", "details": "Valid authentication token is required"})
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
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Document service is unavailable - please try again later",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Not Found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":      "Document not found",
				"details":    "The requested document doesn't exist or you don't have permission to modify it",
				"debug_info": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update document",
			"details": err.Error(),
		})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Document service rejected the update request",
			"details": res.Error,
		})
		return
	}

	// Если есть активные WebSocket соединения для этого документа, уведомляем их об обновлении
	if h.wsService.GetActiveConnections(documentID) > 0 {
		h.wsService.NotifyDocumentUpdated(documentID, userID, res.Document)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"document": res.Document,
	})
}

// DeleteDocument обрабатывает запрос на удаление документа
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	documentID := c.Param("id")
	if documentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing document ID", "details": "Document ID is required in the path"})
		return
	}

	// Получаем ID пользователя из токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required", "details": "Valid authentication token is required"})
		return
	}

	// Проверяем, есть ли активные соединения с этим документом
	hasActiveConnections := h.wsService.GetActiveConnections(documentID) > 0

	// Отправляем запрос к document-сервису через gRPC
	res, err := h.documentClient.DeleteDocument(context.Background(), &pb.DeleteDocumentRequest{
		Id:     documentID,
		UserId: userID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Document service is unavailable - please try again later",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Not Found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":      "Document not found",
				"details":    "The requested document doesn't exist or you don't have permission to delete it",
				"debug_info": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete document",
			"details": err.Error(),
		})
		return
	}

	if !res.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Document service rejected the deletion request",
			"details": res.Error,
		})
		return
	}

	// Если есть активные соединения, отправляем уведомление об удалении документа
	if hasActiveConnections {
		h.wsService.NotifyDocumentDeleted(documentID, userID)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Document successfully deleted",
	})
}
