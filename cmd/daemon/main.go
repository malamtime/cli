package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/malamtime/cli/daemon"
	mc "github.com/malamtime/cli/daemon"
	"github.com/malamtime/cli/model"
)

var (
	version    = "dev"
	commit     = "none"
	date       = "unknown"
	uptraceDsn = ""
)

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home directory", slog.Any("err", err))
		return ""
	}
	return filepath.Join(homeDir, ".shelltime", "config.toml")
}

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(l)
	config := &mc.Config{}

	// Parse command line flags
	flag.StringVar(&config.SocketPath, "socket", mc.DefaultSocketPath, "Unix domain socket path")
	flag.Parse()

	// TODO: read from global config
	cs := model.NewConfigService(getConfigPath())

	daemon.Init(cs)

	pubsub := daemon.NewGoChannel(daemon.PubSubConfig{}, watermill.NewSlogLogger(slog.Default()))
	msg, err := pubsub.Subscribe(context.Background(), daemon.PubSubTopic)

	if err != nil {
		slog.Error("Failed to subscribe the message queue topic", slog.String("topic", daemon.PubSubTopic), slog.Any("err", err))
		return
	}

	go daemon.SocketTopicProccessor(msg)

	// Create processor instance
	processor := daemon.NewSocketHandler(config, pubsub)

	// Start processor
	if err := processor.Start(); err != nil {
		slog.Error("Failed to start processor", slog.Any("err", err))
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Cleanup
	pubsub.Close()
	processor.Stop()
}
