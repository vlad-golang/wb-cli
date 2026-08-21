package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/vlad-golang/wb-cli/common"
)

type Command struct {
	WbClient *common.WbClient
}

func (c *Command) Product() *cli.Command {
	return &cli.Command{
		Name: "product",
		Commands: []*cli.Command{
			{
				Name: "search",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Required: true, Aliases: []string{"n"}},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					resp, err := c.WbClient.WbSearch(ctx, cmd.String("name"))
					if err != nil {
						return fmt.Errorf("search: %w", err)
					}

					text, err := json.MarshalIndent(resp, "", "  ")
					if err != nil {
						return fmt.Errorf("yaml marshal: %w", err)
					}

					fmt.Println(string(text))

					return nil
				},
			},
		},
	}
}
