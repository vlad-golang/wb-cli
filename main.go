package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/philippgille/gokv/bbolt"
	"github.com/urfave/cli/v3"
	"github.com/vlad-golang/wb-cli/command"
	"github.com/vlad-golang/wb-cli/common"
)

type config struct {
	Token          string
	TokenCreatedAt time.Time
}

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

	store, err := bbolt.NewStore(bbolt.Options{Path: filepath.Join(homeDir, "wb-cli", "bbolt.db")})
	if err != nil {
		return fmt.Errorf("open bolt db: %w", err)
	}
	defer store.Close()

	var cfg config
	_, err = store.Get("config", &cfg)
	if err != nil {
		return fmt.Errorf("get token created: %w", err)
	}

	if time.Now().After(cfg.TokenCreatedAt.Add(4 * time.Hour)) {
		cfg.Token, err = getToken(ctx)
		if err != nil {
			return fmt.Errorf("get token: %w", err)
		}

		cfg.TokenCreatedAt = time.Now()

		err = store.Set("config", &cfg)
		if err != nil {
			return fmt.Errorf("save config: %w", err)
		}
	}

	c := command.Command{WbClient: &common.WbClient{Token: cfg.Token}}

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "verbose",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "toon, json. toon - token efficient. json - you can use jq tool",
				Value: "toon",
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
