package document

import (
	"time"

	"github.com/google/uuid"
)

// Document представляет документ пользователя
type Document struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// GetDocumentsRequest представляет запрос на получение списка документов
type GetDocumentsRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

// GetDocumentRequest представляет запрос на получение документа
type GetDocumentRequest struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

// CreateDocumentRequest представляет запрос на создание документа
type CreateDocumentRequest struct {
	Title   string    `json:"title" binding:"required"`
	Content string    `json:"content"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
}

// UpdateDocumentRequest представляет запрос на обновление документа
type UpdateDocumentRequest struct {
	ID      uuid.UUID `json:"id" binding:"required"`
	Title   string    `json:"title" binding:"required"`
	Content string    `json:"content"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
}

// DeleteDocumentRequest представляет запрос на удаление документа
type DeleteDocumentRequest struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
