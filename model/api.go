package model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"error"`
}

type TrackingData struct {
	Shell     string `json:"shell"`
	SessionID int64  `json:"sessionId"`
	Command   string `json:"command"`
	Hostname  string `json:"hostname"`
	Username  string `json:"username"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Result    int    `json:"result"`
	OS        string `json:"os"`
	OSVersion string `json:"osVersion"`
}

type PostTrackArgs struct {
	// nano timestamp
	CursorID int64          `json:"cursorId"`
	Data     []TrackingData `json:"data"`
}

func SendLocalDataToServer(ctx context.Context, config ShellTimeConfig, cursor time.Time, trackingData []TrackingData) error {
	data := PostTrackArgs{
		CursorID: cursor.UnixNano(),
		Data:     trackingData,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", config.APIEndpoint+"/api/v1/track", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("shelltimeCLI@%s", commitID))
	req.Header.Set("Authorization", "CLI "+config.Token)

	logrus.Traceln("http: ", req.URL.String(), len(trackingData))

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	defer resp.Body.Close()

	logrus.Traceln("http: ", resp.Status)
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	var msg errorResponse
	err = json.Unmarshal(buf, &msg)
	if err != nil {
		logrus.Errorln("Failed to parse error response:", err)
		return err
	}
	logrus.Errorln("Error response:", msg.ErrorMessage)
	return errors.New(msg.ErrorMessage)
}
