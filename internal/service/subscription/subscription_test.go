package subscription_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iamstep4ik/TestTaskOzonBank/graph/model"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/service/subscription"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeAndPublish(t *testing.T) {
	svc := subscription.NewSubscriptionService()
	postID := int64(42)
	ch := make(chan *model.Comment, 1)

	svc.Subscribe(postID, ch)

	comment := &model.Comment{
		ID:       1,
		PostID:   postID,
		AuthorID: uuid.New(),
		Content:  "Nice post!",
	}

	svc.Publish(postID, comment)

	select {
	case received := <-ch:
		assert.Equal(t, comment, received)
	case <-time.After(time.Second):
		t.Fatal("timeout: no comment received on channel")
	}
}

func TestUnsubscribe(t *testing.T) {
	svc := subscription.NewSubscriptionService()
	postID := int64(100)
	ch := make(chan *model.Comment, 1)

	svc.Subscribe(postID, ch)
	svc.Unsubscribe(postID, ch)

	comment := &model.Comment{
		ID:       2,
		PostID:   postID,
		AuthorID: uuid.New(),
		Content:  "Should not be received",
	}

	svc.Publish(postID, comment)

	select {
	case msg, ok := <-ch:
		if ok {
			t.Fatalf("expected channel to be closed, but received: %+v", msg)
		}
		// ok == false means the channel was closed, which is expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout: expected closed channel but nothing happened")
	} // expected: channel is closed and ignored
}

func TestPublishToMultipleSubscribers(t *testing.T) {
	svc := subscription.NewSubscriptionService()
	postID := int64(77)

	ch1 := make(chan *model.Comment, 1)
	ch2 := make(chan *model.Comment, 1)

	svc.Subscribe(postID, ch1)
	svc.Subscribe(postID, ch2)

	comment := &model.Comment{
		ID:       3,
		PostID:   postID,
		AuthorID: uuid.New(),
		Content:  "Multicast!",
	}

	svc.Publish(postID, comment)

	select {
	case msg := <-ch1:
		assert.Equal(t, comment, msg)
	case <-time.After(time.Second):
		t.Error("timeout waiting for ch1")
	}

	select {
	case msg := <-ch2:
		assert.Equal(t, comment, msg)
	case <-time.After(time.Second):
		t.Error("timeout waiting for ch2")
	}
}
