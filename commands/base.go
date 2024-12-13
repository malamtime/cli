package commands

import (
	"github.com/malamtime/cli/model"
	"go.opentelemetry.io/otel"
)

var commitID string

var configService model.ConfigService

var commandTracer = otel.Tracer("cli")

func InjectVar(commitId string, cs model.ConfigService) {
	commitID = commitId
	configService = cs
}
