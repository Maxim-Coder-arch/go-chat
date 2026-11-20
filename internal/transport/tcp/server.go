package tcp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"chat-app/internal/application/chat"
)

type Server struct {
	address string
	handler chat.MessageHandler
	hub     chat.ClientRegistry
	logger  *slog.Logger
}

func NewServer(
	address string,
	handler chat.MessageHandler,
	hub chat.ClientRegistry,
	logger *slog.Logger,
) *Server {
	return &Server{
		address: address,
		handler: handler,
		hub:     hub,
		logger:  logger,
	}
}

func (s *Server) Serve(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("start TCP server: %w", err)
	}

	defer listener.Close()

	stop := context.AfterFunc(ctx, func() {
		_ = listener.Close()
	})
	defer stop()

	s.logger.Info("TCP server started", "address", s.address)

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()

		if err != nil {
			if ctx.Err() != nil {
				break
			}

			s.logger.Error("Failed to accept connection", "error", err)
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			client := NewConnection(
				conn,
				s.handler,
				s.hub,
			)

			client.Run(ctx)
		}()
	}

	wg.Wait()

	s.logger.Info("TCP server stopped")

	return nil
}
