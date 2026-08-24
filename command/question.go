package command

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *Command) Question() *cli.Command {
	return &cli.Command{
		Name: "question",
		Commands: []*cli.Command{
			{
				Name: "list",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "root-id", Aliases: []string{"r"}, Required: true},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					resp, err := c.WbClient.Questions(ctx, cmd.Int("root-id"))
					if err != nil {
						return fmt.Errorf("wb client questions: %w", err)
					}

					err = printResponse(&resp, cmd)
					if err != nil {
						return fmt.Errorf("print response: %w", err)
					}

					return nil
				},
			},
		},
	}
}
