package postservice

import (
	"context"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/storage"
	"go.uber.org/zap"
)

type PostService struct {
	storage storage.Storage
	log     *zap.Logger
}

func NewPostService(storage storage.Storage, logger *zap.Logger) *PostService {
	return &PostService{storage: storage, log: logger}
}

func (s *PostService) CreatePost(ctx context.Context, newPost *model.NewPost) (*model.Post, error) {
	s.log.Debug("Creating new post", zap.String("authorID", newPost.AuthorID.String()), zap.String("title", newPost.Title))
	post, err := s.storage.CreatePost(ctx, newPost)
	if err != nil {
		s.log.Error("Failed to create post", zap.Error(err), zap.String("authorID", newPost.AuthorID.String()), zap.String("title", newPost.Title))
		return nil, err
	}
	s.log.Debug("Successfully created post", zap.Int64("postID", post.ID), zap.String("authorID", post.AuthorID.String()), zap.String("title", post.Title))
	return post, nil
}
func (s *PostService) GetPost(ctx context.Context, id int64) (*model.Post, error) {
	s.log.Debug("Fetching post", zap.Int64("postID", id))
	post, err := s.storage.GetPost(ctx, id)
	if err != nil {
		s.log.Error("Failed to get post", zap.Int64("postID", id), zap.Error(err))
		return nil, err
	}
	s.log.Debug("Successfully fetched post", zap.Int64("postID", id), zap.String("authorID", post.AuthorID.String()))
	return post, nil
}
func (s *PostService) GetPosts(ctx context.Context) ([]*model.Post, error) {
	s.log.Debug("Fetching all posts")
	posts, err := s.storage.GetPosts(ctx)
	if err != nil {
		s.log.Error("Failed to get posts", zap.Error(err))
		return nil, err
	}
	s.log.Debug("Successfully fetched posts", zap.Int("count", len(posts)))
	return posts, nil
}
func (s *PostService) AllowComments(ctx context.Context, authorID string, postID int64, allowed bool) (*model.Post, error) {
	s.log.Debug("Allowing comments for post", zap.String("authorID", authorID), zap.Int64("postID", postID), zap.Bool("allowed", allowed))
	post, err := s.storage.AllowComments(ctx, authorID, postID, allowed)
	if err != nil {
		s.log.Error("Failed to allow comments", zap.String("authorID", authorID), zap.Int64("postID", postID), zap.Error(err))
		return nil, err
	}
	s.log.Debug("Successfully allowed comments for post", zap.String("authorID", authorID), zap.Int64("postID", postID), zap.Bool("allowed", allowed))
	return post, nil
}

func (s *PostService) GetCommentsForPost(ctx context.Context, postID int64, offset int64, limit int64) ([]*model.Comment, error) {
	s.log.Debug("Fetching comments for post", zap.Int64("postID", postID), zap.Int64("offset", offset), zap.Int64("limit", limit))
	comments, err := s.storage.GetCommentsForPost(ctx, postID, offset, limit)
	if err != nil {
		s.log.Error("Failed to get comments for post", zap.Int64("postID", postID), zap.Error(err))
		return nil, err
	}
	s.log.Debug("Successfully fetched comments for post", zap.Int64("postID", postID), zap.Int("count", len(comments)))
	return comments, nil
}
