package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
	"github.com/malamtime/cli/model"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/browser"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/trace"
)

var AuthCommand *cli.Command = &cli.Command{
	Name:  "init",
	Usage: "init your shelltime.xyz config",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Authentication token",
			Required: false,
		},
	},
	Action: commandAuth,
}

func commandAuth(c *cli.Context) error {
	ctx, span := commandTracer.Start(c.Context, "auth")
	defer span.End()
	configDir := os.ExpandEnv("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}
	SetupLogger(configDir)

	var config model.ShellTimeConfig
	configFile := configDir + "/config.toml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		content, err := toml.Marshal(model.DefaultConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}
		err = os.WriteFile(configFile, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		config = model.DefaultConfig
	} else {
		existingConfig, err := configService.ReadConfigFile(ctx)
		if err != nil {
			return err
		}
		config = existingConfig
	}

	newToken := c.String("token")

	if newToken == "" {
		nt, err := ApplyTokenByHandshake(ctx, config)
		if err != nil {
			return err
		}
		newToken = nt
	}

	config.Token = newToken
	content, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(configFile, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	color.Green.Println(" ✅ config file created")
	return nil
}

func ApplyTokenByHandshake(_ctx context.Context, config model.ShellTimeConfig) (string, error) {
	ctx, span := commandTracer.Start(_ctx, "handshake", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	hs := model.NewHandshakeService(config)

	hid, err := hs.Init(ctx)
	if err != nil {
		return "", err
	}

	startedAt := time.Now()

	feLink := fmt.Sprintf("%s/cli/integration?hid=%s", config.WebEndpoint, hid)

	if err := browser.OpenURL(feLink); err != nil {
		logrus.Errorln(err)
	}

	color.Green.Println(fmt.Sprintf("Open %s to continue", feLink))

	s := spinner.New(spinner.CharSets[35], 200*time.Millisecond)
	s.Start()
	for {
		if time.Since(startedAt) > 10*time.Minute {
			color.Red.Println(" ❌ Failed to authenticate. Please retry with `shelltime init` or contact shelltime team (annatar.he+shelltime.xyz@gmail.com)")
			s.Stop()
			return "", fmt.Errorf("authentication timeout")
		}

		token, err := hs.Check(ctx, hid)
		if err != nil {
			return "", err
		}
		if token != "" {
			color.Green.Println(" ✅ You are ready to go!")
			s.Stop()
			return token, nil
		}

		time.Sleep(2 * time.Second)
	}
}
