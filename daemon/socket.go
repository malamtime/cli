package daemon

import (
	"log/slog"
	"net"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/vmihailenco/msgpack/v5"
)

type SocketMessageType string

const (
	SocketMessageTypeSync SocketMessageType = "sync"
)

type SocketMessage struct {
	Type    SocketMessageType `msg:"type"`
	Payload interface{}       `msg:"payload"`
}

type SocketHandler struct {
	config   *Config
	listener net.Listener

	channel  *GoChannel
	stopChan chan struct{}
}

func NewSocketHandler(config *Config, ch *GoChannel) *SocketHandler {
	return &SocketHandler{
		config:   config,
		channel:  ch,
		stopChan: make(chan struct{}),
	}
}

func (p *SocketHandler) Start() error {
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

func (p *SocketHandler) Stop() {
	p.channel.Close()
	close(p.stopChan)
	if p.listener != nil {
		p.listener.Close()
	}
	os.RemoveAll(p.config.SocketPath)
	slog.Info("Daemon stopped")
}

func (p *SocketHandler) acceptConnections() {
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

func (p *SocketHandler) handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := msgpack.NewDecoder(conn)
	var msg SocketMessage
	if err := decoder.Decode(&msg); err != nil {
		slog.Error("Error decoding message", slog.Any("err", err))
		return
	}

	switch msg.Type {
	// case "status":
	// 	p.handleStatus(conn)
	// case "track":
	// 	p.handleTrack(conn, msg.Payload)
	case SocketMessageTypeSync:
		buf, err := msgpack.Marshal(msg)
		if err != nil {
			slog.Error("Error encoding message", slog.Any("err", err))
		}

		chMsg := message.NewMessage(watermill.NewUUID(), buf)
		if err := p.channel.Publish(PubSubTopic, chMsg); err != nil {
			slog.Error("Error to publish topic", slog.Any("err", err))
		}
	default:
		slog.Error("Unknown message type:", slog.String("messageType", string(msg.Type)))
	}
}
