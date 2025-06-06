package subscription

import (
	"sync"

	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
)

type SubscriptionService struct {
	Subscribers map[int64][]chan *model.Comment
	mu          sync.Mutex
}

func NewSubscriptionService() *SubscriptionService {
	return &SubscriptionService{
		Subscribers: make(map[int64][]chan *model.Comment),
		mu:          sync.Mutex{},
	}
}

func (s *SubscriptionService) Subscribe(postID int64, channel chan *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Subscribers[postID] = append(s.Subscribers[postID], channel)
}
