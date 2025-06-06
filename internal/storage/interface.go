package storage

import (
	"context"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
)

type StorageMethods interface {
	CreatePost(ctx context.Context, newPost *model.NewPost) (*model.Post, error)
	CreateComment(ctx context.Context, newComment *model.NewComment) (*model.Comment, error)
	AllowComments(ctx context.Context, authorID string, postID int64, allowed bool) (*model.Post, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetPost(ctx context.Context, id int64) (*model.Post, error)
	GetCommentsForPost(ctx context.Context, postID int64, offset int64, limit int64) ([]*model.Comment, error)
}
