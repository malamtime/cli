// cli/daemon/chan_test.go
package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type chanTestSuite struct {
	suite.Suite
}

func (s *chanTestSuite) TestNewGoChannel() {
	s.T().Run("default configuration", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{}, nil)
		assert.NotNil(t, pubSub)
		assert.NotNil(t, pubSub.subscribers)
		assert.NotNil(t, pubSub.persistedMessages)
		assert.NotNil(t, pubSub.closing)
	})
}

func (s *chanTestSuite) TestPublishSubscribe() {
	s.T().Run("simple publish subscribe", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{
			OutputChannelBuffer: 10,
		}, nil)
		defer pubSub.Close()

		msg := message.NewMessage("1", []byte("test"))
		topic := "test-topic"

		// Create subscriber first
		messages, err := pubSub.Subscribe(context.Background(), topic)
		assert.NoError(t, err)

		// Publish message
		err = pubSub.Publish(topic, msg)
		assert.NoError(t, err)

		// Receive message
		received := <-messages
		assert.Equal(t, msg.UUID, received.UUID)
		assert.Equal(t, msg.Payload, received.Payload)
	})

	s.T().Run("persistent messages", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{
			OutputChannelBuffer: 10,
			Persistent:          true,
		}, nil)
		defer pubSub.Close()

		msg := message.NewMessage("1", []byte("test"))
		topic := "test-topic"

		// Publish before subscribe
		err := pubSub.Publish(topic, msg)
		assert.NoError(t, err)

		// Subscribe after publish
		messages, err := pubSub.Subscribe(context.Background(), topic)
		assert.NoError(t, err)

		// Should receive persisted message
		received := <-messages
		assert.Equal(t, msg.UUID, received.UUID)
		assert.Equal(t, msg.Payload, received.Payload)
	})

	s.T().Run("multiple subscribers", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{
			OutputChannelBuffer: 10,
		}, nil)
		defer pubSub.Close()

		msg := message.NewMessage("1", []byte("test"))
		topic := "test-topic"

		// Create two subscribers
		messages1, err := pubSub.Subscribe(context.Background(), topic)
		assert.NoError(t, err)
		messages2, err := pubSub.Subscribe(context.Background(), topic)
		assert.NoError(t, err)

		// Publish message
		err = pubSub.Publish(topic, msg)
		assert.NoError(t, err)

		// Both subscribers should receive the message
		received1 := <-messages1
		received2 := <-messages2
		assert.Equal(t, msg.UUID, received1.UUID)
		assert.Equal(t, msg.UUID, received2.UUID)
	})

	s.T().Run("subscriber context cancellation", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{
			OutputChannelBuffer: 10,
		}, nil)
		defer pubSub.Close()

		ctx, cancel := context.WithCancel(context.Background())
		topic := "test-topic"

		_, err := pubSub.Subscribe(ctx, topic)
		assert.NoError(t, err)

		// Cancel context
		cancel()

		// Wait a bit for cleanup
		time.Sleep(100 * time.Millisecond)

		// Verify subscriber was removed
		assert.Empty(t, pubSub.subscribers[topic])
	})
}

func (s *chanTestSuite) TestClose() {
	s.T().Run("close with active subscribers", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{
			OutputChannelBuffer: 10,
		}, nil)

		topic := "test-topic"
		_, err := pubSub.Subscribe(context.Background(), topic)
		assert.NoError(t, err)

		err = pubSub.Close()
		assert.NoError(t, err)

		// Verify closed state
		assert.True(t, pubSub.isClosed())
		assert.Nil(t, pubSub.persistedMessages)
	})

	s.T().Run("publish to closed channel", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{}, nil)
		pubSub.Close()

		msg := message.NewMessage("1", []byte("test"))
		err := pubSub.Publish("topic", msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Pub/Sub closed")
	})

	s.T().Run("subscribe to closed channel", func(t *testing.T) {
		pubSub := NewGoChannel(PubSubConfig{}, nil)
		pubSub.Close()

		_, err := pubSub.Subscribe(context.Background(), "topic")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Pub/Sub closed")
	})
}

func TestChanTestSuite(t *testing.T) {
	suite.Run(t, new(chanTestSuite))
}
