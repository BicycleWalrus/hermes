package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/BicycleWalrus/hermes/pkg/kubeeye"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := kubeeye.Run(ctx, logger); err != nil {
		logger.Error("kubeeye exited with error", "error", err)
		os.Exit(1)
	}
}
