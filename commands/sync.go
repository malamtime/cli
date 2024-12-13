package commands

import (
	"os"

	"github.com/malamtime/cli/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/trace"
)

var SyncCommand *cli.Command = &cli.Command{
	Name:   "sync",
	Usage:  "manually sync local commands to server",
	Action: commandSync,
	OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
		return nil
	},
}

func commandSync(c *cli.Context) error {
	ctx, span := commandTracer.Start(c.Context, "sync", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	SetupLogger(os.ExpandEnv("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER))

	config, err := configService.ReadConfigFile(ctx)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return trySyncLocalToServer(ctx, config, true)
}
