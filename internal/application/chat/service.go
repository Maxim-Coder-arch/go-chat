
package chat

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "time"

    "chat-app/internal/domain"
)

type Service struct {
    store      MessageStore
    broadcaster Broadcaster
}

func NewService(
    store MessageStore,
    broadcaster Broadcaster,
) *Service {
    return &Service{
        store:       store,
        broadcaster: broadcaster,
    }
}

func (s *Service) SendMessage(
    ctx context.Context,
    username string,
    text string,
) (domain.Message, error) {
    message := domain.Message{
        ID:        generateID(),
        Username:  username,
        Text:      text,
        CreatedAt: time.Now().UTC(),
    }

    if err := message.Validate(); err != nil {
        return domain.Message{}, err
    }

    if err := s.store.Save(ctx, message); err != nil {
        return domain.Message{}, err
    }

    s.broadcaster.Broadcast(message)

    return message, nil
}

func generateID() string {
    id := make([]byte, 16)

    if _, err := rand.Read(id); err != nil {
        panic(err)
    }

    return hex.EncodeToString(id)
}