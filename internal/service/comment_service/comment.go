package commentservice

import (
	"context"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/log"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/storage"
	"go.uber.org/zap"
)

type CommentService struct {
	storage storage.StorageMethods
	log     *zap.Logger
}

func NewCommentService(storage storage.StorageMethods, logger *zap.Logger) *CommentService {
	return &CommentService{storage: storage, log: log.GetLogger()}
}

func (s *CommentService) CreateComment(ctx context.Context, newComment *model.NewComment) (*model.Comment, error) {
	s.log.Info("Creating new comment", zap.Any("newComment", newComment))
	comment, err := s.storage.CreateComment(ctx, newComment)
	if err != nil {
		s.log.Error("Failed to create comment", zap.Error(err))
		return nil, err
	}
	s.log.Info("Comment created successfully", zap.Any("comment", comment))
	return comment, nil

}
