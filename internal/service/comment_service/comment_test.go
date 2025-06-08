package commentservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	commentservice "github.com/iamstep4ik/TestTaskOzonBank/internal/service/comment_service"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/mocks"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateComment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()
	service := commentservice.NewCommentService(mockStorage, logger)

	authorID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	input := &model.NewComment{
		AuthorID: authorID,
		PostID:   1,
		Content:  "Hello",
	}
	expected := &model.Comment{ID: 123, AuthorID: authorID, Content: "Hello"}

	mockStorage.
		EXPECT().
		CreateComment(gomock.Any(), input).
		Return(expected, nil)

	result, err := service.CreateComment(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != expected.ID {
		t.Errorf("expected %v, got %v", expected.ID, result.ID)
	}
}

func TestGetReplies_InvalidLimit(t *testing.T) {
	service := commentservice.NewCommentService(nil, zap.NewNop())
	_, err := service.GetReplies(context.Background(), 1, nil, ptr(int64(101)))
	if err == nil || err.Error() != "maximum limit is 100" {
		t.Errorf("expected max limit error, got: %v", err)
	}
}

func TestGetCommentDepth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	logger := zap.NewNop()
	service := commentservice.NewCommentService(mockStorage, logger)

	mockStorage.
		EXPECT().
		GetCommentDepth(gomock.Any(), int64(42)).
		Return(3, nil)

	depth, err := service.GetCommentDepth(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if depth != 3 {
		t.Errorf("expected depth 3, got %v", depth)
	}
}

func ptr(i int64) *int64 { return &i }
