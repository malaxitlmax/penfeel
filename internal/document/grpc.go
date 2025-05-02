package document

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/malaxitlmax/penfeel/api/proto"
)

// GRPCServer реализация gRPC сервера для документов
type GRPCServer struct {
	pb.DocumentServiceServer
	service Service
}

// NewGRPCServer создает новый gRPC сервер для документов
func NewGRPCServer(service Service) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}

// GetDocuments обрабатывает запрос на получение списка документов
func (s *GRPCServer) GetDocuments(ctx context.Context, req *pb.GetDocumentsRequest) (*pb.GetDocumentsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.GetDocumentsResponse{
			Success: false,
			Error:   "invalid user ID",
		}, nil
	}

	// Преобразуем запрос в доменную модель
	domainReq := GetDocumentsRequest{
		UserID: userID,
	}

	// Вызываем сервис для получения документов
	documents, err := s.service.GetDocuments(ctx, domainReq)
	if err != nil {
		return &pb.GetDocumentsResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Преобразуем документы в protobuf формат
	var pbDocuments []*pb.Document
	for _, doc := range documents {
		pbDocuments = append(pbDocuments, &pb.Document{
			Id:        doc.ID.String(),
			Title:     doc.Title,
			Content:   doc.Content,
			UserId:    doc.UserID.String(),
			CreatedAt: doc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: doc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	// Формируем ответ
	return &pb.GetDocumentsResponse{
		Success:   true,
		Documents: pbDocuments,
	}, nil
}

// GetDocument обрабатывает запрос на получение документа
func (s *GRPCServer) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.GetDocumentResponse{
			Success: false,
			Error:   "invalid document ID",
		}, nil
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.GetDocumentResponse{
			Success: false,
			Error:   "invalid user ID",
		}, nil
	}

	// Преобразуем запрос в доменную модель
	domainReq := GetDocumentRequest{
		ID:     id,
		UserID: userID,
	}

	// Вызываем сервис для получения документа
	document, err := s.service.GetDocument(ctx, domainReq)
	if err != nil {
		return &pb.GetDocumentResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Формируем ответ
	return &pb.GetDocumentResponse{
		Success: true,
		Document: &pb.Document{
			Id:        document.ID.String(),
			Title:     document.Title,
			Content:   document.Content,
			UserId:    document.UserID.String(),
			CreatedAt: document.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: document.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// CreateDocument обрабатывает запрос на создание документа
func (s *GRPCServer) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*pb.CreateDocumentResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.CreateDocumentResponse{
			Success: false,
			Error:   "invalid user ID",
		}, nil
	}

	// Преобразуем запрос в доменную модель
	domainReq := CreateDocumentRequest{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	// Вызываем сервис для создания документа
	document, err := s.service.CreateDocument(ctx, domainReq)
	if err != nil {
		return &pb.CreateDocumentResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Формируем ответ
	return &pb.CreateDocumentResponse{
		Success: true,
		Document: &pb.Document{
			Id:        document.ID.String(),
			Title:     document.Title,
			Content:   document.Content,
			UserId:    document.UserID.String(),
			CreatedAt: document.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: document.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// UpdateDocument обрабатывает запрос на обновление документа
func (s *GRPCServer) UpdateDocument(ctx context.Context, req *pb.UpdateDocumentRequest) (*pb.UpdateDocumentResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.UpdateDocumentResponse{
			Success: false,
			Error:   "invalid document ID",
		}, nil
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.UpdateDocumentResponse{
			Success: false,
			Error:   "invalid user ID",
		}, nil
	}

	// Преобразуем запрос в доменную модель
	domainReq := UpdateDocumentRequest{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	// Вызываем сервис для обновления документа
	document, err := s.service.UpdateDocument(ctx, domainReq)
	if err != nil {
		return &pb.UpdateDocumentResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Формируем ответ
	return &pb.UpdateDocumentResponse{
		Success: true,
		Document: &pb.Document{
			Id:        document.ID.String(),
			Title:     document.Title,
			Content:   document.Content,
			UserId:    document.UserID.String(),
			CreatedAt: document.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: document.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// DeleteDocument обрабатывает запрос на удаление документа
func (s *GRPCServer) DeleteDocument(ctx context.Context, req *pb.DeleteDocumentRequest) (*pb.DeleteDocumentResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return &pb.DeleteDocumentResponse{
			Success: false,
			Error:   "invalid document ID",
		}, nil
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.DeleteDocumentResponse{
			Success: false,
			Error:   "invalid user ID",
		}, nil
	}

	// Преобразуем запрос в доменную модель
	domainReq := DeleteDocumentRequest{
		ID:     id,
		UserID: userID,
	}

	// Вызываем сервис для удаления документа
	err = s.service.DeleteDocument(ctx, domainReq)
	if err != nil {
		return &pb.DeleteDocumentResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	// Формируем ответ
	return &pb.DeleteDocumentResponse{
		Success: true,
	}, nil
}
