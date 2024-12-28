package daemon

import "github.com/malamtime/cli/model"

var stConfig model.ConfigService

const (
	PubSubTopic = "socket"
)

func Init(cs model.ConfigService) {
	stConfig = cs
}
