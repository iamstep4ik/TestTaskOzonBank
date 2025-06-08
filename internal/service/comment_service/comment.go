package commentservice

import (
	"context"
	"fmt"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/storage"
	"go.uber.org/zap"
)

type CommentService struct {
	storage storage.Storage
	log     *zap.Logger
}

func NewCommentService(storage storage.Storage, logger *zap.Logger) *CommentService {
	return &CommentService{storage: storage, log: logger}
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

func (s *CommentService) GetReplies(ctx context.Context, commentID int64, offset, limit *int64) ([]*model.Comment, error) {

	off := int64(0)
	if offset != nil {
		off = *offset
	}

	lim := int64(10)
	if limit != nil {
		lim = *limit
	}

	if lim > 100 {
		return nil, fmt.Errorf("maximum limit is 100")
	}
	if off < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}

	s.log.Debug("Fetching comment replies",
		zap.Int64("comment_id", commentID),
		zap.Int64("offset", off),
		zap.Int64("limit", lim))

	replies, err := s.storage.GetRepliesByParentID(ctx, commentID, off, lim)
	if err != nil {
		s.log.Error("Failed to get comment replies",
			zap.Error(err),
			zap.Int64("comment_id", commentID))
		return nil, fmt.Errorf("failed to get replies: %w", err)
	}

	s.log.Info("Successfully fetched comment replies",
		zap.Int64("comment_id", commentID),
		zap.Int("reply_count", len(replies)))
	return replies, nil
}

func (s *CommentService) GetCommentDepth(ctx context.Context, commentID int64) (int, error) {
	s.log.Debug("Calculating comment depth",
		zap.Int64("comment_id", commentID))

	depth, err := s.storage.GetCommentDepth(ctx, commentID)
	if err != nil {
		s.log.Error("Failed to calculate comment depth",
			zap.Error(err),
			zap.Int64("comment_id", commentID))
		return 0, fmt.Errorf("failed to calculate depth: %w", err)
	}

	s.log.Info("Comment depth calculated",
		zap.Int64("comment_id", commentID),
		zap.Int("depth", depth))
	return depth, nil
}
