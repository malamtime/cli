package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
}

type PostTrackArgs struct {
	Data []TrackingData `json:"data"`
}

func SendLocalDataToServer(ctx context.Context, config MalamTimeConfig, trackingData []TrackingData) error {
	data := PostTrackArgs{
		Data: trackingData,
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
	req.Header.Set("User-Agent", fmt.Sprintf("MalamTimeCLI@%s", commitID))
	req.Header.Set("X-API", "Bearer "+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		logrus.Errorln(resp.Status)
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorln(err)
		}
		var msg errorResponse
		if err := json.Unmarshal(buf, &msg); err != nil {
			logrus.Errorln("Failed to parse error response:", err)
		} else {
			logrus.Errorln("Error response:", msg.ErrorMessage)
		}
		return err
	}
	return nil
}
