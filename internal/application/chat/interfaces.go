package chat

import (
	"chat-app/internal/domain"
	"context"
)

type MessageStore interface {
	Save(ctx context.Context, message domain.Message) error
	GetRecent(ctx context.Context, limit int) ([]domain.Message, error)
}

type Broadcaster interface {
	Broadcast(message domain.Message)
}

type Client interface {
    Send(message domain.Message) error
    Close() error
}

type MessageHandler interface {
    SendMessage(
        ctx context.Context,
        username string,
        text string,
    ) (domain.Message, error)
}