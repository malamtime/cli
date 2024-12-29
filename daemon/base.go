package daemon

import "github.com/malamtime/cli/model"

var stConfig model.ConfigService
var version string

const (
	PubSubTopic = "socket"
)

func Init(cs model.ConfigService, vs string) {
	stConfig = cs
	version = vs
}
