package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gocarina/gocsv"
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
					&cli.StringFlag{Name: "query", Required: true, Aliases: []string{"q"}},
					&cli.IntFlag{Name: "page", Aliases: []string{"p"}, Value: 1, Usage: "results page number"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					resp, err := c.WbClient.WbSearch(ctx, cmd.String("query"), cmd.Int("page"))
					if err != nil {
						return fmt.Errorf("search: %w", err)
					}

					err = gocsv.Marshal(&resp, cmd.Writer)
					if err != nil {
						return fmt.Errorf("csv marshal: %w", err)
					}

					return nil
				},
			},
			{
				Name: "get",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "id", Aliases: []string{"i"}, Required: true},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					resp, err := c.WbClient.GetCard(ctx, cmd.Int("id"))
					if err != nil {
						return fmt.Errorf("wb client get card: %w", err)
					}

					err = gocsv.Marshal(&resp, cmd.Writer)
					if err != nil {
						err := json.NewEncoder(cmd.Writer).Encode(&resp)
						if err != nil {
							return fmt.Errorf("json marshal: %w", err)
						}
					}

					return nil
				},
			},
		},
	}
}
