package main

import (
	"context"

	"github.com/chromedp/chromedp"
)

func NewContext(ctx context.Context) (context.Context, context.CancelFunc) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-site-isolation-trials", true),
		// chromedp.Flag("disable-web-security", true), // Breaks sites
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		// chromedp.UserDataDir("tmp/chromedp"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		// /Users/vmsharoshkin/Library/Application Support/Google/Chrome
		// chromedp.Flag("profile-directory", "Chromedp Profile"),
	)

	cancelFunctions := make([]context.CancelFunc, 0, 2)

	ctx, cancel := chromedp.NewExecAllocator(ctx, opts...)

	cancelFunctions = append(cancelFunctions, cancel)

	// if cfg.Allocator == "remote" {
	//	ctx, cancel = chromedp.NewRemoteAllocator(
	//		context.Background(),
	//		"ws://localhost:9222", // Try https://lightpanda.io
	//	)
	//	cancelFuncs = append(cancelFuncs, cancel)
	//}

	ctx, cancel = chromedp.NewContext(ctx) // chromedp.WithDebugf(log.Printf))
	cancelFunctions = append(cancelFunctions, cancel)

	return ctx, func() {
		for _, c := range cancelFunctions {
			c()
		}
	}
}
