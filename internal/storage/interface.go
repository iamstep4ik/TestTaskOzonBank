package storage

import (
	"context"
	"os"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	dbs "github.com/iamstep4ik/TestTaskOzonBank/internal/storage/db"
	inmemory "github.com/iamstep4ik/TestTaskOzonBank/internal/storage/in-memory"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	CreatePost(ctx context.Context, newPost *model.NewPost) (*model.Post, error)
	CreateComment(ctx context.Context, newComment *model.NewComment) (*model.Comment, error)
	AllowComments(ctx context.Context, authorID string, postID int64, allowed bool) (*model.Post, error)
	GetPosts(ctx context.Context) ([]*model.Post, error)
	GetPost(ctx context.Context, id int64) (*model.Post, error)
	GetCommentsForPost(ctx context.Context, postID int64, offset int64, limit int64) ([]*model.Comment, error)
	GetCommentDepth(ctx context.Context, commentID int64) (int, error)
	GetRepliesByParentID(ctx context.Context, parentID int64, offset, limit int64) ([]*model.Comment, error)
}

type StorageType string

const (
	StorageTypeDB     StorageType = "db"
	StorageTypeMemory StorageType = "memory"
)

func NewStorage(ctx context.Context, db *pgxpool.Pool) Storage {
	storageType := StorageType(os.Getenv("STORAGE_TYPE"))

	switch storageType {
	case StorageTypeMemory:
		return inmemory.NewStorageMemory()
	case StorageTypeDB:
		if db == nil {
			panic("database connection is required for database storage")
		}
		return dbs.NewStorageDB(db)
	default:
		if db != nil {
			return dbs.NewStorageDB(db)
		}
		return inmemory.NewStorageMemory()
	}
}
