package inmemory

import (
	"context"
	"time"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/utils/errs"
)

type StorageMemory struct {
	posts          map[int64]*model.Post
	comments       map[int64][]*model.Comment
	commentMap     map[int64]*model.Comment
	postCounter    int64
	commentCounter int64
}

func NewStorageMemory() *StorageMemory {
	return &StorageMemory{
		posts:          make(map[int64]*model.Post),
		comments:       make(map[int64][]*model.Comment),
		commentMap:     make(map[int64]*model.Comment),
		postCounter:    0,
		commentCounter: 0,
	}
}

func (s *StorageMemory) CreatePost(ctx context.Context, newPost *model.NewPost) (*model.Post, error) {
	post := &model.Post{
		ID:              s.postCounter,
		AuthorID:        newPost.AuthorID,
		Title:           newPost.Title,
		Content:         newPost.Content,
		CommentsAllowed: newPost.CommentsAllowed,
		CreatedAt:       time.Now(),
	}
	s.posts[post.ID] = post
	s.postCounter++
	return post, nil
}

func (s *StorageMemory) CreateComment(ctx context.Context, newComment *model.NewComment) (*model.Comment, error) {
	post, exists := s.posts[newComment.PostID]
	if !exists {
		return nil, errs.ErrPostNotFound
	}
	if !post.CommentsAllowed {
		return nil, errs.ErrCommentsNotAllowed
	}

	if newComment.ParentID != nil {
		if _, exists := s.commentMap[*newComment.ParentID]; !exists {
			return nil, errs.ErrParentCommentNotFound
		}
	}

	comment := &model.Comment{
		ID:        s.commentCounter,
		AuthorID:  newComment.AuthorID,
		PostID:    newComment.PostID,
		ParentID:  newComment.ParentID,
		Content:   newComment.Content,
		CreatedAt: time.Now(),
	}

	s.comments[comment.PostID] = append(s.comments[comment.PostID], comment)
	s.commentMap[comment.ID] = comment
	s.commentCounter++

	return comment, nil
}

func (s *StorageMemory) AllowComments(ctx context.Context, authorID string, postID int64, allowed bool) (*model.Post, error) {
	post, exists := s.posts[postID]
	if !exists || post.AuthorID.String() != authorID {
		return nil, errs.ErrPostNotFound
	}
	post.CommentsAllowed = allowed
	return post, nil
}

func (s *StorageMemory) GetPosts(ctx context.Context) ([]*model.Post, error) {
	posts := make([]*model.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *StorageMemory) GetPost(ctx context.Context, id int64) (*model.Post, error) {
	if post, exists := s.posts[id]; exists {
		return post, nil
	}
	return nil, errs.ErrPostNotFound
}

func (s *StorageMemory) GetCommentsForPost(ctx context.Context, postID int64, offset int64, limit int64) ([]*model.Comment, error) {
	comments, exists := s.comments[postID]
	if !exists {
		return nil, nil
	}

	if offset >= int64(len(comments)) {
		return []*model.Comment{}, nil
	}

	end := offset + limit
	if end > int64(len(comments)) {
		end = int64(len(comments))
	}

	return comments[offset:end], nil
}

func (s *StorageMemory) GetRepliesByParentID(ctx context.Context, parentID int64, offset, limit int64) ([]*model.Comment, error) {
	var replies []*model.Comment
	for _, comments := range s.comments {
		for _, comment := range comments {
			if comment.ParentID != nil && *comment.ParentID == parentID {
				replies = append(replies, comment)
			}
		}
	}

	if offset >= int64(len(replies)) {
		return []*model.Comment{}, nil
	}

	end := offset + limit
	if end > int64(len(replies)) {
		end = int64(len(replies))
	}

	return replies[offset:end], nil
}

func (s *StorageMemory) GetCommentDepth(ctx context.Context, commentID int64) (int, error) {
	depth := 0
	currentID := commentID

	for {
		comment, exists := s.commentMap[currentID]
		if !exists || comment.ParentID == nil {
			break
		}
		depth++
		currentID = *comment.ParentID
	}

	return depth, nil
}
