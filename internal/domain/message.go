package domain

import (
	"strings"
	"time"
)

type Message struct {
	ID        string
	Username  string
	Text      string
	CreatedAt time.Time
}

func (m Message) Validate() error {
	if strings.TrimSpace(m.Username) == "" {
		return ErrEmptyUsername
	}

	if strings.TrimSpace(m.Text) == "" {
		return ErrEmptyMessage
	}

	return nil
}
