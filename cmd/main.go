package main

import (
	"context"
	"fmt"
	cli "github.com/artarts36/singlecli"
	"github.com/artarts36/telegram-action/internal"
	"github.com/caarlos0/env/v11"
)

type config struct {
	Token           string `env:"INPUT_TOKEN,required"`
	ChatID          string `env:"INPUT_CHAT_ID,required"`
	ChatThreadID    string `env:"INPUT_CHAT_THREAD_ID"`
	Message         string `env:"INPUT_MESSAGE,required"`
	Host            string `env:"INPUT_HOST,required"`
	ConvertMarkdown bool   `env:"INPUT_CONVERT_MARKDOWN"`
	IssueTrackerURL string `env:"INPUT_ISSUE_TRACKER_URL"`
}

func main() {
	app := &cli.App{
		BuildInfo: &cli.BuildInfo{
			Name: "telegram-action",
		},
		Args:   []*cli.ArgDefinition{},
		Action: run,
	}

	app.RunWithGlobalArgs(context.Background())
}

const colorGreen = 2

func run(ctx *cli.Context) error {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	msg := internal.Message{
		Body:            cfg.Message,
		ChatID:          cfg.ChatID,
		ChatThreadID:    cfg.ChatThreadID,
		ConvertMarkdown: cfg.ConvertMarkdown,
	}

	if cfg.IssueTrackerURL != "" {
		tracker := internal.NewIssueTracker(cfg.IssueTrackerURL)
		msg.Body = tracker.InjectLinks(msg.Body)
	}

	messenger := internal.NewMessenger(cfg.Token, cfg.Host, internal.NewMarkdownToHTMLConverter())

	res, err := messenger.Send(ctx.Context, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	ctx.Output.PrintColoredBlock(colorGreen, fmt.Sprintf("Message was sent. ID: %d", res.MessageID))

	return nil
}
