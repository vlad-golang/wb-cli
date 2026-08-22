package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

					searchResp := make([]SearchResponse, 0, len(resp.Products))
					for _, product := range resp.Products {
						searchResp = append(searchResp, SearchResponse{
							ID:             product.ID,
							Name:           product.Name,
							Brand:          product.Brand,
							PriceRub:       product.Sizes[0].Price.Product / 100,
							ReviewRating:   product.ReviewRating,
							SupplierRating: product.SupplierRating,
							Feedbacks:      product.Feedbacks,
							RootID:         int(product.Root),
						})
					}

					err = printResponse(&searchResp)
					if err != nil {
						return fmt.Errorf("print response: %w", err)
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

					out := ProductGetOut{
						NmID:        resp.NmID,
						ImtID:       resp.ImtID,
						Name:        resp.ImtName,
						Category:    resp.SubjName,
						Brand:       resp.Selling.BrandName,
						Description: resp.Description,
						Options:     resp.Options,
						Verbose:     resp,
					}

					if !cmd.IsSet("verbose") {
						out.Verbose = common.GetCardBody{}
					}

					err = printResponse(&out)
					if err != nil {
						return fmt.Errorf("print response: %w", err)
					}

					return nil
				},
			},
		},
	}
}

func printResponse(v any) error {
	jsonString, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return fmt.Errorf("json marshal options: %w", err)
	}

	fmt.Println(string(jsonString))

	return nil
}

type SearchResponse struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Brand          string  `json:"brand"`
	PriceRub       int     `json:"price_rub"`
	ReviewRating   float64 `json:"review_rating"`
	SupplierRating float64 `json:"supplier_rating"`
	Feedbacks      int     `json:"feedbacks"`
	RootID         int     `json:"root_id"`
}

type ProductGetOut struct {
	NmID        int                `json:"nm_id"`
	ImtID       int                `json:"imt_id"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Category    string             `json:"category"`
	VendorCode  string             `json:"vendor_code"`
	Brand       string             `json:"brand"`
	SupplierID  int                `json:"supplier_id"`
	Description string             `json:"description"`
	Options     []common.Options   `json:"options"`
	Colors      []int              `json:"colors"`
	PhotoCount  int                `json:"photo_count"`
	HasVideo    bool               `json:"has_video"`
	CreatedAt   time.Time          `json:"created_at"`
	Verbose     common.GetCardBody `json:"verbose"`
}
