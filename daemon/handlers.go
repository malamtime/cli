package daemon

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/vmihailenco/msgpack/v5"
)

func SocketTopicProccessor(messages <-chan *message.Message) {
	for msg := range messages {
		ctx := context.Background()
		fmt.Printf("received message: %s, payload: %s\n", msg.UUID, string(msg.Payload))

		var socketMsg SocketMessage
		if err := msgpack.Unmarshal(msg.Payload, socketMsg); err != nil {
			slog.ErrorContext(ctx, "failed to parse socket message", slog.Any("err", err))
			return
		}

		if socketMsg.Type == "sync" {
			if err := handlePubSubSync(ctx, socketMsg.Payload); err != nil {
				slog.ErrorContext(ctx, "failed to parse socket message", slog.Any("err", err))
				return
			}
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
