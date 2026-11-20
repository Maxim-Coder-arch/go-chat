package memory

import (
	"context"
	"sync"

	"chat-app/internal/domain"
)

const DefaultHistoryLimit = 100

type MessageStore struct {
	mu       sync.RWMutex
	messages []domain.Message
	limit    int
}

func NewMessageStore(limit int) *MessageStore {
	if limit <= 0 {
		limit = DefaultHistoryLimit
	}

	return &MessageStore{
		messages: make([]domain.Message, 0, limit),
		limit:    limit,
	}
}

func (s *MessageStore) Save(ctx context.Context, message domain.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	if len(s.messages) >= s.limit {
		copy(s.messages, s.messages[1:])
		s.messages = s.messages[:s.limit-1]
	}

	s.messages = append(s.messages, message)

	return nil
}

func (s *MessageStore) GetRecent(ctx context.Context, limit int) ([]domain.Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if limit <= 0 {
		return []domain.Message{}, nil
	}

	if limit > len(s.messages) {
		limit = len(s.messages)
	}

	start := len(s.messages) - limit

	result := make([]domain.Message, limit)
	copy(result, s.messages[start:])

	return result, nil
}

var _ interface {
	Save(context.Context, domain.Message) error
	GetRecent(context.Context, int) ([]domain.Message, error)
} = (*MessageStore)(nil)