package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vlad-golang/wb-cli/command"
	"github.com/vlad-golang/wb-cli/common"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/urfave/cli/v3"
)

type config struct {
	Token          string
	TokenCreatedAt time.Time
}

// Мелкие неудобства: описания команд в help пустые, product get без форматирования JSON, флага --version нет.
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("os user home dir: %w", err)
	}

	configFile := filepath.Join(homeDir, "wb-cli", "config.json")

	cfg := config{}
	err = cleanenv.ReadConfig(configFile, &cfg)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if time.Now().After(cfg.TokenCreatedAt.Add(23 * time.Hour)) {
		cfg.Token, err = getToken(ctx)
		if err != nil {
			return fmt.Errorf("get token: %w", err)
		}

		cfg.TokenCreatedAt = time.Now()

		marshaledCfg, err := json.Marshal(&cfg)
		if err != nil {
			return fmt.Errorf("marshal config: %w", err)
		}

		err = os.WriteFile(configFile, marshaledCfg, 0o600)
		if err != nil {
			return fmt.Errorf("write config: %w", err)
		}
	}

	c := command.Command{WbClient: &common.WbClient{Token: cfg.Token}}

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose output",
			},
		},
		Commands: []*cli.Command{
			c.Product(),
			c.Feedback(),
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

func getToken(ctx context.Context) (string, error) {
	tokenChan := make(chan string, 1)

	ctx, cancel := NewContext(ctx)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev any) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
		case *network.EventResponseReceived:
		case *network.EventRequestWillBeSentExtraInfo:
			if cookie, ok := ev.Headers["cookie"]; ok {
				tokenChan <- cookie.(string)
				cancel()
			}

		}
	})

	err := chromedp.Run(ctx, network.Enable(), chromedp.Navigate("https://wildberries.ru"))
	if err != nil && !errors.Is(err, context.Canceled) {
		return "", fmt.Errorf("error getting token: %w", err)
	}

	return <-tokenChan, nil
}
