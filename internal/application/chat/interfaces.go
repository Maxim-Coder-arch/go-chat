package chat

import (
	"context"
	"chat-app/internal/domain"
)

type MessageStore interface {
	Save(ctx context.Context, message domain.Message) error
    GetRecent(ctx context.Context, limit int) ([]domain.Message, error)
}

type Broadcaster interface {
    Broadcast(message domain.Message)
}