package daemon

const (
	DefaultSocketPath = "/tmp/shelltime.sock"
)

type Config struct {
	SocketPath string
}
