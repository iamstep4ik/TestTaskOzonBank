package mapper

import (
	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/models"
)

func PostToGen(dbPost *models.Post) *model.Post {
	return &model.Post{
		ID:              dbPost.ID,
		AuthorID:        dbPost.AuthorID,
		Title:           dbPost.Title,
		Content:         dbPost.Content,
		CommentsAllowed: dbPost.CommentsAllowed,
		CreateDate:      dbPost.CreateDate,
	}
}
