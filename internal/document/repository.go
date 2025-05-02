package document

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository интерфейс репозитория для работы с документами
type Repository interface {
	GetDocuments(ctx context.Context, userID uuid.UUID) ([]*Document, error)
	GetDocument(ctx context.Context, id, userID uuid.UUID) (*Document, error)
	CreateDocument(ctx context.Context, doc *Document) (*Document, error)
	UpdateDocument(ctx context.Context, doc *Document) (*Document, error)
	DeleteDocument(ctx context.Context, id, userID uuid.UUID) error
}

// PostgresRepository реализация репозитория для PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository создает новый репозиторий для PostgreSQL
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// GetDocuments возвращает список документов пользователя
func (r *PostgresRepository) GetDocuments(ctx context.Context, userID uuid.UUID) ([]*Document, error) {
	var documents []*Document
	query := `SELECT * FROM documents WHERE user_id = $1 ORDER BY updated_at DESC`
	err := r.db.SelectContext(ctx, &documents, query, userID)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

// GetDocument возвращает документ по ID
func (r *PostgresRepository) GetDocument(ctx context.Context, id, userID uuid.UUID) (*Document, error) {
	var document Document
	query := `SELECT * FROM documents WHERE id = $1 AND user_id = $2`
	err := r.db.GetContext(ctx, &document, query, id, userID)
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// CreateDocument создает новый документ
func (r *PostgresRepository) CreateDocument(ctx context.Context, doc *Document) (*Document, error) {
	query := `INSERT INTO documents (title, content, user_id) 
              VALUES ($1, $2, $3) 
              RETURNING id, title, content, user_id, created_at, updated_at`

	var document Document
	err := r.db.QueryRowxContext(ctx, query, doc.Title, doc.Content, doc.UserID).
		StructScan(&document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// UpdateDocument обновляет документ
func (r *PostgresRepository) UpdateDocument(ctx context.Context, doc *Document) (*Document, error) {
	query := `UPDATE documents 
              SET title = $1, content = $2, updated_at = $3
              WHERE id = $4 AND user_id = $5
              RETURNING id, title, content, user_id, created_at, updated_at`

	now := time.Now()
	var document Document
	err := r.db.QueryRowxContext(ctx, query, doc.Title, doc.Content, now, doc.ID, doc.UserID).
		StructScan(&document)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// DeleteDocument удаляет документ
func (r *PostgresRepository) DeleteDocument(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM documents WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}
