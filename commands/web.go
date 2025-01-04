// cli/commands/web.go
package commands

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/trace"
)

var WebCommand *cli.Command = &cli.Command{
	Name:   "web",
	Usage:  "open ShellTime web dashboard in browser",
	Action: commandWeb,
	OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
		color.Red.Println(err.Error())
		return nil
	},
}

func commandWeb(c *cli.Context) error {
	ctx, span := commandTracer.Start(c.Context, "web", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	config, err := configService.ReadConfigFile(ctx)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	if config.WebEndpoint == "" {
		return fmt.Errorf("web endpoint is not configured. Please check your configuration")
	}

	url := config.WebEndpoint
	err = browser.OpenURL(url)
	if err != nil {
		return fmt.Errorf("failed to open browser: %v", err)
	}

	color.Green.Printf("Opening %s in your default browser\n", url)
	return nil
}
