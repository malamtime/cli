package daemon

import (
	"context"
	"log/slog"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/pkg/errors"
)

func handlePubSubSync(ctx context.Context, socketMsgPayload interface{}) error {
	syncMsg, ok := socketMsgPayload.(model.PostTrackArgs)
	if !ok {
		slog.Error("Failed to parse sync payload", slog.Any("payload", socketMsgPayload))
		return errors.New("failed to parse the payload")
	}

	cfg, err := stConfig.ReadConfigFile(ctx)
	if err != nil {
		slog.Error("Failed to unmarshal sync message", slog.Any("err", err))
		return err
	}

	// Call SendLocalDataToServer
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
