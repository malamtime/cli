package model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"error"`
}

type TrackingData struct {
	Shell     string `json:"shell" msgpack:"shell"`
	SessionID int64  `json:"sessionId" msgpack:"sessionId"`
	Command   string `json:"command" msgpack:"command"`
	Hostname  string `json:"hostname" msgpack:"hostname"`
	Username  string `json:"username" msgpack:"username"`
	StartTime int64  `json:"startTime" msgpack:"startTime"`
	EndTime   int64  `json:"endTime" msgpack:"endTime"`
	Result    int    `json:"result" msgpack:"result"`
	OS        string `json:"os" msgpack:"os"`
	OSVersion string `json:"osVersion" msgpack:"osVersion"`
}

type PostTrackArgs struct {
	// nano timestamp
	CursorID int64          `json:"cursorId" msgpack:"cursorId"`
	Data     []TrackingData `json:"data" msgpack:"data"`
}

func doSendData(ctx context.Context, endpoint Endpoint, cursor time.Time, trackingData []TrackingData) error {
	data := PostTrackArgs{
		CursorID: cursor.UnixNano(),
		Data:     trackingData,
	}

	jsonData, err := msgpack.Marshal(data)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 3,
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint.APIEndpoint+"/api/v1/track", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	req.Header.Set("Content-Type", "application/msgpack")
	req.Header.Set("User-Agent", fmt.Sprintf("shelltimeCLI@%s", commitID))
	req.Header.Set("Authorization", "CLI "+endpoint.Token)

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

func SendLocalDataToServer(ctx context.Context, config ShellTimeConfig, cursor time.Time, trackingData []TrackingData) error {
	if config.Token == "" {
		logrus.Traceln("no token available. do not sync to server")
		return nil
	}

	var wg sync.WaitGroup

	wg.Add(len(config.Endpoints) + 1)

	authPair := make([]Endpoint, len(config.Endpoints)+1)

	authPair[0] = Endpoint{
		Token:       config.Token,
		APIEndpoint: config.APIEndpoint,
	}

	copy(authPair[1:], config.Endpoints)

	errs := make(chan error, len(authPair))

	for _, pair := range authPair {
		go func(pair Endpoint) {
			defer wg.Done()
			err := doSendData(ctx, pair, cursor, trackingData)
			errs <- err
		}(pair)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
