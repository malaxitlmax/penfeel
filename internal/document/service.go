package document

import (
	"context"
	"fmt"
)

// Service интерфейс сервиса для работы с документами
type Service interface {
	GetDocuments(ctx context.Context, req GetDocumentsRequest) ([]*Document, error)
	GetDocument(ctx context.Context, req GetDocumentRequest) (*Document, error)
	CreateDocument(ctx context.Context, req CreateDocumentRequest) (*Document, error)
	UpdateDocument(ctx context.Context, req UpdateDocumentRequest) (*Document, error)
	DeleteDocument(ctx context.Context, req DeleteDocumentRequest) error
}

// DocumentService реализация сервиса для работы с документами
type DocumentService struct {
	repo Repository
}

// NewDocumentService создает новый сервис для работы с документами
func NewDocumentService(repo Repository) *DocumentService {
	return &DocumentService{
		repo: repo,
	}
}

// GetDocuments возвращает список документов пользователя
func (s *DocumentService) GetDocuments(ctx context.Context, req GetDocumentsRequest) ([]*Document, error) {
	documents, err := s.repo.GetDocuments(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	return documents, nil
}

// GetDocument возвращает документ по ID
func (s *DocumentService) GetDocument(ctx context.Context, req GetDocumentRequest) (*Document, error) {
	document, err := s.repo.GetDocument(ctx, req.ID, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return document, nil
}

// CreateDocument создает новый документ
func (s *DocumentService) CreateDocument(ctx context.Context, req CreateDocumentRequest) (*Document, error) {
	document := &Document{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	}

	createdDoc, err := s.repo.CreateDocument(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return createdDoc, nil
}

// UpdateDocument обновляет документ
func (s *DocumentService) UpdateDocument(ctx context.Context, req UpdateDocumentRequest) (*Document, error) {
	document := &Document{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	}

	updatedDoc, err := s.repo.UpdateDocument(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return updatedDoc, nil
}

// DeleteDocument удаляет документ
func (s *DocumentService) DeleteDocument(ctx context.Context, req DeleteDocumentRequest) error {
	err := s.repo.DeleteDocument(ctx, req.ID, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
