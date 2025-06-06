package db

import (
	"context"
	"time"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/log"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/utils/errs"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		db:  db,
		log: log.GetLogger(),
	}
}

func (r *Storage) CreatePost(ctx context.Context, newPost *model.NewPost) (*model.Post, error) {
	r.log.Info("Creating new post", zap.String("author_id", newPost.AuthorID.String()))
	post := &model.Post{
		AuthorID:        newPost.AuthorID,
		Title:           newPost.Title,
		Content:         newPost.Content,
		CommentsAllowed: newPost.CommentsAllowed,
		CreatedAt:       time.Now(),
	}

	query := `INSERT INTO posts (author_id, title, content, allow_comments, created_at)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING post_id`
	err := r.db.QueryRow(ctx, query, post.AuthorID, post.Title, post.Content, post.CommentsAllowed, post.CreatedAt).Scan(&post.ID)
	if err != nil {
		r.log.Error("Failed to create post", zap.Error(err), zap.String("author_id", post.AuthorID.String()))
		return nil, err
	}

	r.log.Info("Post created", zap.Int64("post_id", post.ID), zap.String("author_id", post.AuthorID.String()))

	return post, nil
}

func (r *Storage) CreateComment(ctx context.Context, newComment *model.NewComment) (*model.Comment, error) {
	r.log.Info("Creating new comment", zap.String("author_id", newComment.AuthorID.String()), zap.Int64("post_id", newComment.PostID))

	queryCommentsAllowed := `SELECT allow_comments FROM posts WHERE post_id = $1`
	var commentsAllowed bool
	err := r.db.QueryRow(ctx, queryCommentsAllowed, newComment.PostID).Scan(&commentsAllowed)
	if err != nil {
		r.log.Error("Failed to check if comments are allowed", zap.Error(err), zap.Int64("post_id", newComment.PostID))
		return nil, err
	}

	if !commentsAllowed {
		r.log.Warn("Comments are not allowed for this post", zap.Int64("post_id", newComment.PostID))
		return nil, errs.ErrCommentsNotAllowed
	}

	if newComment.ParentID != nil {
		var parentExists bool
		err := r.db.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM comments WHERE comment_id = $1 AND post_id = $2)`,
			*newComment.ParentID,
			newComment.PostID,
		).Scan(&parentExists)

		if err != nil {
			r.log.Error("Failed to check if parent comment exists", zap.Error(err), zap.Int64("parent_id", *newComment.ParentID), zap.Int64("post_id", newComment.PostID))
			return nil, err
		}
		if !parentExists {
			r.log.Warn("Parent comment does not exist", zap.Int64("parent_id", *newComment.ParentID), zap.Int64("post_id", newComment.PostID))
			return nil, err
		}
	}

	queryCommentCreate := `
        INSERT INTO comments (author_id, post_id, parent_id, content, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING comment_id, created_at
    `

	comment := &model.Comment{
		AuthorID:  newComment.AuthorID,
		PostID:    newComment.PostID,
		ParentID:  newComment.ParentID,
		Content:   newComment.Content,
		CreatedAt: time.Now(),
	}

	err = r.db.QueryRow(
		ctx,
		queryCommentCreate,
		comment.AuthorID,
		comment.PostID,
		comment.ParentID,
		comment.Content,
		comment.CreatedAt,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		r.log.Error("Failed to create comment", zap.Error(err), zap.String("author_id", comment.AuthorID.String()), zap.Int64("post_id", comment.PostID))
		return nil, err
	}
	r.log.Info("Comment created", zap.Int64("comment_id", comment.ID), zap.String("author_id", comment.AuthorID.String()), zap.Int64("post_id", comment.PostID))

	return comment, nil
}

func (r *Storage) AllowComments(ctx context.Context, authorID string, postID int64, allowed bool) (*model.Post, error) {
	r.log.Info("Updating comments allowed for post", zap.Int64("post_id", postID), zap.String("author_id", authorID), zap.Bool("allowed", allowed))

	query := `UPDATE posts SET allow_comments = $1 WHERE post_id = $2 AND author_id = $3 RETURNING post_id, author_id, title, content, allow_comments, created_at`
	post := &model.Post{}
	err := r.db.QueryRow(ctx, query, allowed, postID, authorID).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CommentsAllowed, &post.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.log.Warn("Post not found or author mismatch", zap.Int64("post_id", postID), zap.String("author_id", authorID))
			return nil, errs.ErrPostNotFound
		}
		r.log.Error("Failed to update comments allowed", zap.Error(err), zap.Int64("post_id", postID), zap.String("author_id", authorID))
		return nil, err
	}

	r.log.Info("Comments allowed updated", zap.Int64("post_id", post.ID), zap.String("author_id", post.AuthorID.String()), zap.Bool("allowed", post.CommentsAllowed))
	return post, nil
}
func (r *Storage) GetPosts(ctx context.Context) ([]*model.Post, error) {
	r.log.Info("Fetching all posts")
	query := `SELECT post_id, author_id, title, content, allow_comments, created_at
			  FROM posts`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.log.Error("Failed to fetch posts", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CommentsAllowed, &post.CreatedAt)
		if err != nil {
			r.log.Error("Failed to scan post", zap.Error(err))
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		r.log.Error("Failed to fetch posts", zap.Error(err))
		return nil, err
	}
	r.log.Info("Posts fetched successfully", zap.Int("count", len(posts)))
	return posts, nil
}
func (r *Storage) GetPost(ctx context.Context, id int64) (*model.Post, error) {
	query := `SELECT post_id, author_id, title, content, allow_comments, created_at
			  FROM posts WHERE post_id = $1`
	post := &model.Post{}
	err := r.db.QueryRow(ctx, query, id).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CommentsAllowed, &post.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			r.log.Warn("Post not found", zap.Int64("post_id", id))
			return nil, errs.ErrPostNotFound
		}
		r.log.Error("Failed to fetch post", zap.Error(err), zap.Int64("post_id", id))
		return nil, err
	}

	return post, nil

}

func (r *Storage) GetCommentsForPost(ctx context.Context, postID int64, offset int64, limit int64) ([]*model.Comment, error) {
	query := `SELECT comment_id, author_id, post_id, parent_id, content, created_at
			  FROM comments WHERE post_id = $1 ORDER BY created_at ASC OFFSET $2 LIMIT $3`
	rows, err := r.db.Query(ctx, query, postID, offset, limit)
	if err != nil {
		r.log.Error("Failed to fetch comments for post", zap.Error(err), zap.Int64("post_id", postID))
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		comment := &model.Comment{}
		err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			r.log.Error("Failed to scan comment", zap.Error(err), zap.Int64("post_id", postID))
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		r.log.Error("Failed to fetch comments for post", zap.Error(err), zap.Int64("post_id", postID))
		return nil, err
	}
	r.log.Info("Comments fetched successfully", zap.Int("count", len(comments)), zap.Int64("post_id", postID))
	return comments, nil
}
