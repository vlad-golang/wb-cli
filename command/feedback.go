package command

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *Command) Feedback() *cli.Command {
	return &cli.Command{
		Name: "feedback",
		Commands: []*cli.Command{
			{
				Name: "list",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "root-id", Required: true, Aliases: []string{"r"}},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					resp, err := c.WbClient.Feedbacks(ctx, cmd.Int("root-id"))
					if err != nil {
						return fmt.Errorf("wb client feedbacks: %w", err)
					}

					err = printResponse(&resp)
					if err != nil {
						return fmt.Errorf("print response: %w", err)
					}

					return nil
				},
			},
		},
	}
}
