package mapper

import (
	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/models"
)

func PostToDB(post *model.NewPost) *models.NewPost {
	return &models.NewPost{
		AuthorID:        post.AuthorID,
		Title:           post.Title,
		Content:         post.Content,
		CommentsAllowed: post.CommentsAllowed,
	}
}
