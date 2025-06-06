package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        int64     `json:"id" db:"comment_id"`
	AuthorID  uuid.UUID `json:"authorID" db:"author_id"`
	PostID    int64     `json:"postID" db:"post_id"`
	ParentID  *int64    `json:"parentID,omitempty" db:"parent_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type NewComment struct {
	AuthorID uuid.UUID `json:"authorID"  db:"author_id"`
	PostID   int64     `json:"postID" db:"post_id"`
	ParentID *int64    `json:"parentID,omitempty" db:"parent_id"`
	Content  string    `json:"content" db:"content"`
}

type NewPost struct {
	AuthorID        uuid.UUID `json:"authorID"  db:"author_id"`
	Title           string    `json:"title" db:"title"`
	Content         string    `json:"content" db:"content"`
	CommentsAllowed bool      `json:"commentsAllowed" db:"comments_allowed"`
}

type Post struct {
	ID              int64      `json:"id" db:"post_id"`
	AuthorID        uuid.UUID  `json:"authorID" db:"author_id"`
	Title           string     `json:"title" db:"title"`
	Content         string     `json:"content" db:"content"`
	CommentsAllowed bool       `json:"commentsAllowed" db:"comments_allowed"`
	Comments        []*Comment `json:"comments,omitempty"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
}
