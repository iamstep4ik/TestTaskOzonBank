package postservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/mocks"
	postservice "github.com/iamstep4ik/TestTaskOzonBank/internal/service/post_service"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreatePost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()
	service := postservice.NewPostService(mockStorage, logger)

	authorID := uuid.New()
	input := &model.NewPost{AuthorID: authorID, Title: "Test", Content: "Content"}
	expected := &model.Post{ID: 1, AuthorID: authorID, Title: "Test", Content: "Content"}

	mockStorage.EXPECT().
		CreatePost(gomock.Any(), input).
		Return(expected, nil)

	result, err := service.CreatePost(context.Background(), input)
	if err != nil || result.ID != expected.ID {
		t.Errorf("CreatePost failed: got %v, expected %v", result, expected)
	}
}

func TestGetPost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()
	service := postservice.NewPostService(mockStorage, logger)

	postID := int64(1)
	expected := &model.Post{ID: postID, AuthorID: uuid.New(), Title: "Post", Content: "..."}

	mockStorage.EXPECT().
		GetPost(gomock.Any(), postID).
		Return(expected, nil)

	result, err := service.GetPost(context.Background(), postID)
	if err != nil || result.ID != postID {
		t.Errorf("GetPost failed: got %v, expected %v", result, expected)
	}
}

func TestGetPosts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()
	service := postservice.NewPostService(mockStorage, logger)

	expected := []*model.Post{
		{ID: 1, AuthorID: uuid.New(), Title: "One"},
		{ID: 2, AuthorID: uuid.New(), Title: "Two"},
	}

	mockStorage.EXPECT().
		GetPosts(gomock.Any()).
		Return(expected, nil)

	posts, err := service.GetPosts(context.Background())
	if err != nil || len(posts) != 2 {
		t.Errorf("GetPosts failed: got %d posts, expected 2", len(posts))
	}
}

func TestAllowComments_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()
	service := postservice.NewPostService(mockStorage, logger)

	postID := int64(42)
	authorID := uuid.New().String()
	allowed := true

	expected := &model.Post{ID: postID, AuthorID: uuid.MustParse(authorID), Title: "Commented"}

	mockStorage.EXPECT().
		AllowComments(gomock.Any(), authorID, postID, allowed).
		Return(expected, nil)

	post, err := service.AllowComments(context.Background(), authorID, postID, allowed)
	if err != nil || post.ID != postID {
		t.Errorf("AllowComments failed: got %v, expected post ID %d", post, postID)
	}
}

func TestGetCommentsForPost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()
	service := postservice.NewPostService(mockStorage, logger)

	postID := int64(1)
	offset, limit := int64(0), int64(10)
	expected := []*model.Comment{
		{ID: 1, PostID: postID, Content: "Nice!"},
		{ID: 2, PostID: postID, Content: "Cool"},
	}

	mockStorage.EXPECT().
		GetCommentsForPost(gomock.Any(), postID, offset, limit).
		Return(expected, nil)

	comments, err := service.GetCommentsForPost(context.Background(), postID, offset, limit)
	if err != nil || len(comments) != 2 {
		t.Errorf("GetCommentsForPost failed: got %v, expected %v", comments, expected)
	}
}
