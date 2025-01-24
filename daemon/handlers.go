package daemon

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/vmihailenco/msgpack/v5"
)

func SocketTopicProccessor(messages <-chan *message.Message) {
	for msg := range messages {
		ctx := context.Background()
		slog.InfoContext(ctx, "received message: ", slog.String("msg.uuid", msg.UUID))

		var socketMsg SocketMessage
		if err := msgpack.Unmarshal(msg.Payload, &socketMsg); err != nil {
			slog.ErrorContext(ctx, "failed to parse socket message", slog.Any("err", err))
		}

		if socketMsg.Type == SocketMessageTypeSync {
			if err := handlePubSubSync(ctx, socketMsg.Payload); err != nil {
				slog.ErrorContext(ctx, "failed to parse socket message", slog.Any("err", err))
			}
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
