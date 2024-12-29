package daemon

import (
	"context"
	"log/slog"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/vmihailenco/msgpack/v5"
)

func handlePubSubSync(ctx context.Context, socketMsgPayload interface{}) error {
	pb, err := msgpack.Marshal(socketMsgPayload)
	if err != nil {
		slog.Error("Failed to marshal the sync payload again for unmarshal", slog.Any("payload", socketMsgPayload))
		return err
	}

	var syncMsg model.PostTrackArgs
	err = msgpack.Unmarshal(pb, &syncMsg)
	if err != nil {
		slog.Error("Failed to parse sync payload", slog.Any("payload", socketMsgPayload))
		return err
	}

	cfg, err := stConfig.ReadConfigFile(ctx)
	if err != nil {
		slog.Error("Failed to unmarshal sync message", slog.Any("err", err))
		return err
	}

	// Call SendLocalDataToServer
	slog.Debug("Sending local data to server",
		slog.Any("ctx", ctx),
		slog.Any("cfg", cfg),
		slog.Time("cursor", time.Unix(0, syncMsg.CursorID)),
		slog.Any("data", syncMsg.Data),
		slog.Any("meta", syncMsg.Meta),
	)

	err = model.SendLocalDataToServer(
		ctx,
		cfg,
		time.Unix(0, syncMsg.CursorID), // Convert nano timestamp to time.Time
		syncMsg.Data,
		syncMsg.Meta,
	)

	if err != nil {
		slog.Error("Failed to sync data to server", slog.Any("err", err))
		return err
	}
	return nil
}
