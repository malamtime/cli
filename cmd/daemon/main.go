package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	mc "github.com/malamtime/cli/daemon"
	"github.com/malamtime/cli/handlers"
	"github.com/malamtime/cli/model"
)

var (
	version    = "dev"
	commit     = "none"
	date       = "unknown"
	uptraceDsn = ""
)

func main() {
	config := &mc.Config{}

	// Parse command line flags
	flag.StringVar(&config.SocketPath, "socket", mc.DefaultSocketPath, "Unix domain socket path")
	flag.Parse()

	// TODO: read from global config
	cs := model.NewConfigService("/home/annatarhe/.config/shelltime/config.toml")

	handlers.Init(cs)

	// Create processor instance
	processor := handlers.NewProcessor(config)

	// Start processor
	if err := processor.Start(); err != nil {
		slog.Error("Failed to start processor", slog.Any("err", err))
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Cleanup
	processor.Stop()
}
