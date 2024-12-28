// cli/daemon/handlers_test.go
package daemon

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/malamtime/cli/model"
	"github.com/malamtime/cli/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vmihailenco/msgpack/v5"
)

type handlersTestSuite struct {
	suite.Suite
}

func (s *handlersTestSuite) TestSocketTopicProcessor() {
	mockedST := mocks.NewConfigService(s.T())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	mockedST.On("ReadConfigFile", mock.Anything).Return(model.ShellTimeConfig{
		APIEndpoint: ts.URL,
	})

	stConfig = mockedST

	s.T().Run("valid sync message", func(t *testing.T) {
		// Create a message channel
		msgChan := make(chan *message.Message)

		// Create a sync message
		socketMsg := SocketMessage{
			Type:    "sync",
			Payload: []byte(`{"some":"data"}`),
		}
		payload, err := msgpack.Marshal(socketMsg)
		assert.NoError(t, err)

		msg := message.NewMessage("test-uuid", payload)

		// Start processor in a goroutine
		go SocketTopicProccessor(msgChan)

		// Send message
		msgChan <- msg

		// Give some time for processing
		time.Sleep(100 * time.Millisecond)

		// Close channel
		close(msgChan)
	})

	s.T().Run("invalid message format", func(t *testing.T) {
		msgChan := make(chan *message.Message)

		// Create an invalid message
		msg := message.NewMessage("test-uuid", []byte("invalid"))

		go SocketTopicProccessor(msgChan)

		// Send invalid message
		msgChan <- msg

		// Give some time for processing
		time.Sleep(100 * time.Millisecond)

		close(msgChan)
	})

	s.T().Run("non-sync message type", func(t *testing.T) {
		msgChan := make(chan *message.Message)

		// Create a non-sync message
		socketMsg := SocketMessage{
			Type:    "other",
			Payload: []byte(`{"some":"data"}`),
		}
		payload, err := msgpack.Marshal(socketMsg)
		assert.NoError(t, err)

		msg := message.NewMessage("test-uuid", payload)

		go SocketTopicProccessor(msgChan)

		// Send message
		msgChan <- msg

		// Give some time for processing
		time.Sleep(100 * time.Millisecond)

		close(msgChan)
	})

	s.T().Run("invalid sync payload", func(t *testing.T) {
		msgChan := make(chan *message.Message)

		// Create a sync message with invalid payload
		socketMsg := SocketMessage{
			Type:    "sync",
			Payload: []byte(`invalid json`),
		}
		payload, err := msgpack.Marshal(socketMsg)
		assert.NoError(t, err)

		msg := message.NewMessage("test-uuid", payload)

		go SocketTopicProccessor(msgChan)

		// Send message
		msgChan <- msg

		// Give some time for processing
		time.Sleep(100 * time.Millisecond)

		close(msgChan)
	})

	s.T().Run("multiple messages", func(t *testing.T) {
		msgChan := make(chan *message.Message)

		// Create multiple messages
		socketMsg1 := SocketMessage{
			Type:    "sync",
			Payload: []byte(`{"some":"data1"}`),
		}
		payload1, err := msgpack.Marshal(socketMsg1)
		assert.NoError(t, err)

		socketMsg2 := SocketMessage{
			Type:    "sync",
			Payload: []byte(`{"some":"data2"}`),
		}
		payload2, err := msgpack.Marshal(socketMsg2)
		assert.NoError(t, err)

		msg1 := message.NewMessage("test-uuid-1", payload1)
		msg2 := message.NewMessage("test-uuid-2", payload2)

		go SocketTopicProccessor(msgChan)

		// Send multiple messages
		msgChan <- msg1
		msgChan <- msg2

		// Give some time for processing
		time.Sleep(100 * time.Millisecond)

		close(msgChan)
	})
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}
