package service

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	pb "github.com/malaxitlmax/penfeel/api/proto"
	"golang.org/x/net/context"
)

// WebSocketService управляет WebSocket соединениями и обработкой сообщений
type WebSocketService struct {
	documentClient pb.DocumentServiceClient

	// Мьютекс для безопасного доступа к карте соединений
	connectionsLock sync.RWMutex
	// documentID -> map[userID]*websocket.Conn
	documentConnections map[string]map[string]*websocket.Conn
}

// NewWebSocketService создаёт новый сервис для обработки WebSocket соединений
func NewWebSocketService(documentClient pb.DocumentServiceClient) *WebSocketService {
	return &WebSocketService{
		documentClient:      documentClient,
		documentConnections: make(map[string]map[string]*websocket.Conn),
	}
}

// RegisterConnection регистрирует новое WebSocket соединение для документа
func (s *WebSocketService) RegisterConnection(documentID, userID string, conn *websocket.Conn) []string {
	s.connectionsLock.Lock()
	defer s.connectionsLock.Unlock()

	// Создаём карту соединений для документа, если её ещё нет
	if _, exists := s.documentConnections[documentID]; !exists {
		s.documentConnections[documentID] = make(map[string]*websocket.Conn)
	}

	// Регистрируем соединение
	s.documentConnections[documentID][userID] = conn

	// Получаем список активных пользователей
	activeUsers := make([]string, 0, len(s.documentConnections[documentID]))
	for uid := range s.documentConnections[documentID] {
		activeUsers = append(activeUsers, uid)
	}

	log.Printf("User %s connected to document %s. Total active users: %d", userID, documentID, len(activeUsers))
	return activeUsers
}

// RemoveConnection удаляет соединение пользователя
func (s *WebSocketService) RemoveConnection(documentID, userID string) {
	s.connectionsLock.Lock()
	defer s.connectionsLock.Unlock()

	// Проверяем, существует ли карта для документа
	if connections, exists := s.documentConnections[documentID]; exists {
		// Удаляем соединение пользователя
		delete(connections, userID)

		// Если соединений для документа больше нет, удаляем карту документа
		if len(connections) == 0 {
			delete(s.documentConnections, documentID)
		}

		log.Printf("User %s disconnected from document %s", userID, documentID)
	}
}

// BroadcastToOthers отправляет сообщение всем пользователям документа, кроме отправителя
func (s *WebSocketService) BroadcastToOthers(documentID string, senderID string, message interface{}) {
	s.connectionsLock.RLock()
	defer s.connectionsLock.RUnlock()

	if connections, exists := s.documentConnections[documentID]; exists {
		for userID, conn := range connections {
			// Не отправляем сообщение отправителю
			if userID != senderID {
				if err := conn.WriteJSON(message); err != nil {
					log.Printf("Error broadcasting to user %s: %v", userID, err)
				}
			}
		}
	}
}

// BroadcastToAll отправляет сообщение всем пользователям документа, включая отправителя
func (s *WebSocketService) BroadcastToAll(documentID string, message interface{}) {
	s.connectionsLock.RLock()
	defer s.connectionsLock.RUnlock()

	if connections, exists := s.documentConnections[documentID]; exists {
		for userID, conn := range connections {
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Error broadcasting to user %s: %v", userID, err)
			}
		}
	}
}

// GetActiveConnections возвращает количество активных соединений для документа
func (s *WebSocketService) GetActiveConnections(documentID string) int {
	s.connectionsLock.RLock()
	defer s.connectionsLock.RUnlock()

	if connections, exists := s.documentConnections[documentID]; exists {
		return len(connections)
	}
	return 0
}

// CloseAllDocumentConnections закрывает все соединения для документа
func (s *WebSocketService) CloseAllDocumentConnections(documentID string) {
	s.connectionsLock.Lock()
	defer s.connectionsLock.Unlock()

	if connections, exists := s.documentConnections[documentID]; exists {
		for _, conn := range connections {
			conn.Close()
		}
		delete(s.documentConnections, documentID)
	}
}

// NotifyDocumentDeleted уведомляет всех пользователей об удалении документа и закрывает соединения
func (s *WebSocketService) NotifyDocumentDeleted(documentID, userID string) {
	// Отправляем уведомление об удалении
	message := map[string]interface{}{
		"type":    "document_deleted",
		"user_id": userID,
	}
	s.BroadcastToAll(documentID, message)

	// Закрываем все соединения
	s.CloseAllDocumentConnections(documentID)
}

// NotifyDocumentUpdated уведомляет всех пользователей об обновлении документа через REST API
func (s *WebSocketService) NotifyDocumentUpdated(documentID, userID string, document *pb.Document) {
	message := map[string]interface{}{
		"type":     "document_updated_externally",
		"document": document,
		"user_id":  userID,
	}
	s.BroadcastToAll(documentID, message)
}

// HandleWebSocketConnection обрабатывает WebSocket соединение после его установки
func (s *WebSocketService) HandleWebSocketConnection(documentID, userID string, conn *websocket.Conn, document *pb.Document) {
	// Получаем список активных пользователей
	activeUsers := s.RegisterConnection(documentID, userID, conn)

	// Отправляем начальное состояние документа
	initialMessage := map[string]interface{}{
		"type":         "init",
		"document":     document,
		"active_users": activeUsers,
	}

	if err := conn.WriteJSON(initialMessage); err != nil {
		log.Println("Error sending initial document state:", err)
		s.RemoveConnection(documentID, userID)
		conn.Close()
		return
	}

	// Оповещаем других пользователей о новом участнике
	s.BroadcastToOthers(documentID, userID, map[string]interface{}{
		"type":    "user_joined",
		"user_id": userID,
	})

	// Устанавливаем отложенное действие для очистки соединения
	defer func() {
		conn.Close()
		s.RemoveConnection(documentID, userID)

		// Оповещаем других пользователей, что пользователь покинул документ
		s.BroadcastToOthers(documentID, userID, map[string]interface{}{
			"type":    "user_left",
			"user_id": userID,
		})
	}()

	// Основной цикл обработки сообщений
	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Обрабатываем сообщение
		s.handleMessage(documentID, userID, conn, rawMessage)
	}
}

// handleMessage обрабатывает входящее WebSocket сообщение
func (s *WebSocketService) handleMessage(documentID, userID string, conn *websocket.Conn, rawMessage []byte) {
	// Декодируем сообщение
	var message map[string]interface{}
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		log.Printf("Error parsing WebSocket message: %v", err)
		return
	}

	messageType, ok := message["type"].(string)
	if !ok {
		log.Println("WebSocket message missing 'type' field")
		return
	}

	log.Printf("Received WebSocket message type: %s from user %s", messageType, userID)

	// Обрабатываем сообщение в зависимости от типа
	switch messageType {
	case "document_update":
		s.handleDocumentUpdate(documentID, userID, conn, message)

	case "cursor_position":
		// Трансляция позиции курсора другим пользователям
		s.BroadcastToOthers(documentID, userID, message)

	case "selection":
		// Трансляция выделения текста другим пользователям
		s.BroadcastToOthers(documentID, userID, message)

	case "ping":
		// Отвечаем на пинг для проверки соединения
		pongMessage := map[string]interface{}{
			"type": "pong",
		}
		if err := conn.WriteJSON(pongMessage); err != nil {
			log.Printf("Error sending pong: %v", err)
		}

	default:
		log.Printf("Unknown message type: %s", messageType)
	}
}

// handleDocumentUpdate обрабатывает обновление документа через WebSocket
func (s *WebSocketService) handleDocumentUpdate(documentID, userID string, conn *websocket.Conn, message map[string]interface{}) {
	content, contentOk := message["content"].(string)
	title, titleOk := message["title"].(string)

	if !contentOk || !titleOk {
		log.Println("Invalid document_update message format")
		return
	}

	// Отправляем изменения в document-сервис
	updateRes, err := s.documentClient.UpdateDocument(context.Background(), &pb.UpdateDocumentRequest{
		Id:      documentID,
		UserId:  userID,
		Title:   title,
		Content: content,
	})

	if err != nil {
		log.Printf("Error updating document: %v", err)
		// Отправляем ошибку только отправителю
		errorMsg := map[string]interface{}{
			"type":  "error",
			"error": "Failed to save document: " + err.Error(),
		}
		conn.WriteJSON(errorMsg)
		return
	}

	if !updateRes.Success {
		log.Printf("Document service rejected update: %s", updateRes.Error)
		errorMsg := map[string]interface{}{
			"type":  "error",
			"error": "Failed to save document: " + updateRes.Error,
		}
		conn.WriteJSON(errorMsg)
		return
	}

	// Если обновление успешно, транслируем изменения другим пользователям
	s.BroadcastToOthers(documentID, userID, message)
}
