package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/malamtime/cli/commands/internal"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var TrackCommand *cli.Command = &cli.Command{
	Name:  "track",
	Usage: "track user commands",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "shell",
			Aliases: []string{"s"},
			Value:   "",
			Usage:   "the shell that user use",
		},
		&cli.Int64Flag{
			Name:    "sessionId",
			Aliases: []string{"id"},
			Value:   0,
			Usage:   "unix timestamp of the session",
		},
		&cli.StringFlag{
			Name:    "command",
			Aliases: []string{"cmd"},
			Value:   "",
			Usage:   "command that user executed",
		},
		&cli.StringFlag{
			Name:    "phase",
			Aliases: []string{"p"},
			Usage:   "Phase: pre, post",
		},
	},
	Action: commandTrack,
	OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
		return nil
	},
}

func commandTrack(c *cli.Context) error {
	logrus.Info(c.Args().First())
	config, err := internal.ReadConfigFile()
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	username := os.Getenv("USER")

	data := map[string]interface{}{
		"shell":     c.String("shell"),
		"sessionId": c.Int64("sessionId"),
		"command":   c.String("command"),
		"hostname":  hostname,
		"username":  username,
		"time":      time.Now().Unix(),
		"phase":     c.String("phase"),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	client := &http.Client{}
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*3)
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
