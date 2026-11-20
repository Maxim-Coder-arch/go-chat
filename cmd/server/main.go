package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"chat-app/internal/application/chat"
	"chat-app/internal/config"
	"chat-app/internal/storage/memory"
	"chat-app/internal/transport/tcp"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	store := memory.NewMessageStore(cfg.HistoryLimit)

	hub := chat.NewHub()

	service := chat.NewService(store, hub)

	server := tcp.NewServer(
		cfg.ServerAddress,
		service,
		hub,
		logger,
	)

	logger.Info(
		"starting chat server",
		"address", cfg.ServerAddress,
		"history_limit", cfg.HistoryLimit,
	)

	if err := server.Serve(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	logger.Info("chat server stopped gracefully")

	return nil
}