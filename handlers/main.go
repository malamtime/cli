package handlers

import (
	"log/slog"
	"net"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	mc "github.com/malamtime/cli/daemon"
	"github.com/vmihailenco/msgpack/v5"
)

type Message struct {
	Type    string      `msg:"type"`
	Payload interface{} `msg:"payload"`
}

type Processor struct {
	config   *mc.Config
	listener net.Listener

	channel  *GoChannel
	stopChan chan struct{}
}

func NewProcessor(config *mc.Config) *Processor {
	ch := NewGoChannel(pubSubConfig{}, watermill.NewSlogLogger(slog.Default()))
	return &Processor{
		config:   config,
		channel:  ch,
		stopChan: make(chan struct{}),
	}
}

func (p *Processor) Start() error {
	// Remove existing socket file if it exists
	if err := os.RemoveAll(p.config.SocketPath); err != nil {
		return err
	}

	// Create Unix domain socket
	listener, err := net.Listen("unix", p.config.SocketPath)
	if err != nil {
		return err
	}
	p.listener = listener

	// Start accepting connections
	go p.acceptConnections()

	slog.Info("Daemon started, listening on: ", slog.String("socketPath", p.config.SocketPath))
	return nil
}

func (p *Processor) Stop() {
	p.channel.Close()
	close(p.stopChan)
	if p.listener != nil {
		p.listener.Close()
	}
	os.RemoveAll(p.config.SocketPath)
	slog.Info("Daemon stopped")
}

func (p *Processor) acceptConnections() {
	for {
		select {
		case <-p.stopChan:
			return
		default:
			conn, err := p.listener.Accept()
			if err != nil {
				continue
			}
			go p.handleConnection(conn)
		}
	}
}

func (p *Processor) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := msgpack.NewDecoder(conn)
	var msg Message
	if err := decoder.Decode(&msg); err != nil {
		slog.Error("Error decoding message", slog.Any("err", err))
		return
	}

	switch msg.Type {
	case "status":
		p.handleStatus(conn)
	case "track":
		p.handleTrack(conn, msg.Payload)
	case "sync":
		p.ProcessSyncMessage(conn, msg.Payload)
	default:
		slog.Error("Unknown message type:", slog.String("messageType", msg.Type))
	}
}

func (p *Processor) handleStatus(conn net.Conn) {
	// Implement status handling
	status := map[string]interface{}{
		"status": "running",
		"uptime": "implement me",
	}
	msgpack.NewEncoder(conn).Encode(status)
}

func (p *Processor) handleTrack(conn net.Conn, payload interface{}) {
	// Implement track handling
	// Save payload to local storage
	response := map[string]interface{}{
		"status": "tracked",
	}
	msgpack.NewEncoder(conn).Encode(response)
}
