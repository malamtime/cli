package handlers

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/vmihailenco/msgpack/v5"
)

// SyncMessage represents the expected structure of a sync message payload
type SyncMessage struct {
	Type         string                 `msgpack:"type"`
	Cursor       int64                  `msgpack:"cursor"`
	TrackingData []model.TrackingData   `msgpack:"trackingData"`
	Meta         model.TrackingMetaData `msgpack:"meta"`
}

// IsSyncMessage checks if the given message is a sync type message
func IsSyncMessage(msg Message) bool {
	return msg.Type == "sync"
}

// ProcessSyncMessage handles the sync message processing
func (p *Processor) ProcessSyncMessage(conn net.Conn, payload interface{}) {
	ctx := context.Background()
	// Create decoder for the payload
	var syncMsg SyncMessage
	payloadBytes, err := msgpack.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal sync payload", slog.Any("err", err))
		return
	}

	err = msgpack.Unmarshal(payloadBytes, &syncMsg)
	if err != nil {
		slog.Error("Failed to unmarshal sync message", slog.Any("err", err))
		return
	}

	cfg, err := stConfig.ReadConfigFile(ctx)
	if err != nil {
		slog.Error("Failed to unmarshal sync message", slog.Any("err", err))
		return
	}

	// Call SendLocalDataToServer
	err = model.SendLocalDataToServer(
		ctx,
		cfg,
		time.Unix(0, syncMsg.Cursor), // Convert nano timestamp to time.Time
		syncMsg.TrackingData,
		syncMsg.Meta,
	)

	if err != nil {
		slog.Error("Failed to sync data to server", slog.Any("err", err))
		return
	}

	// Send success response
	response := map[string]interface{}{
		"status": "synced",
		"cursor": syncMsg.Cursor,
	}
	err = msgpack.NewEncoder(conn).Encode(response)
	if err != nil {
		slog.Error("Failed to send sync response", slog.Any("err", err))
	}
}
