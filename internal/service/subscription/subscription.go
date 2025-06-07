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

func (s *SubscriptionService) Publish(postID int64, comment *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channels, ok := s.Subscribers[postID]; ok {
		var activeChannels []chan *model.Comment

		for _, ch := range channels {
			select {
			case ch <- comment:
				activeChannels = append(activeChannels, ch)
			default:
				close(ch)
			}
		}

		s.Subscribers[postID] = activeChannels
	}
}

func (s *SubscriptionService) Unsubscribe(postID int64, ch chan *model.Comment) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channels, ok := s.Subscribers[postID]; ok {
		for i, channel := range channels {
			if channel == ch {
				close(ch)
				s.Subscribers[postID] = append(channels[:i], channels[i+1:]...)
				break
			}
		}

		if len(s.Subscribers[postID]) == 0 {
			delete(s.Subscribers, postID)
		}
	}
}
