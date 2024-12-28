package daemon

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/vmihailenco/msgpack/v5"
)

func IsSocketReady(ctx context.Context, socketPath string) bool {
	_, err := os.Stat(socketPath)
	return err == nil
}

func SendLocalDataToSocket(ctx context.Context, socketPath string, config model.ShellTimeConfig, cursor time.Time, trackingData []model.TrackingData, meta model.TrackingMetaData) error {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	data := SocketMessage{
		Type: SocketMessageTypeSync,
		Payload: model.PostTrackArgs{
			CursorID: cursor.UnixNano(),
			Data:     trackingData,
			Meta:     meta,
		},
	}

	encoded, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Write(encoded)
	if err != nil {
		return err
	}

	return nil
}
