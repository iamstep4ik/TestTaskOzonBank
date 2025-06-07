package graph

import (
	commentservice "github.com/iamstep4ik/TestTaskOzonBank/internal/service/comment_service"
	postservice "github.com/iamstep4ik/TestTaskOzonBank/internal/service/post_service"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/service/subscription"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService         *postservice.PostService
	CommentService      *commentservice.CommentService
	SubscriptionService *subscription.SubscriptionService
}

func NewResolver(postService *postservice.PostService, commentService *commentservice.CommentService) *Resolver {
	return &Resolver{
		PostService:         postService,
		CommentService:      commentService,
		SubscriptionService: subscription.NewSubscriptionService(),
	}
}
