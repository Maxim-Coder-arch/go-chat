package protocol

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidMessageType = errors.New("invalid message type")

type MessageType string

const (
	MessageTypeJoin    MessageType = "join"
	MessageTypeMessage MessageType = "message"
	MessageTypeLeave   MessageType = "leave"
	MessageTypeSystem  MessageType = "system"
	MessageTypeError   MessageType = "error"
)

type ClientMessage struct {
	Type     MessageType `json:"type"`
	Text     string      `json:"text,omitempty"`
	Username string      `json:"username,omitempty"`
}

func (m ClientMessage) Validate() error {
	switch m.Type {
	case MessageTypeJoin:
		if strings.TrimSpace(m.Username) == "" {
			return errors.New("username cannot be empty")
		}

	case MessageTypeMessage:
		if strings.TrimSpace(m.Text) == "" {
			return errors.New("message cannot be empty")
		}

	case MessageTypeLeave:
		return nil

	default:
		return ErrInvalidMessageType
	}

	return nil
}

type ServerMessage struct {
	Type      MessageType `json:"type"`
	ID        string      `json:"id,omitempty"`
	Username  string      `json:"username,omitempty"`
	Text      string      `json:"text,omitempty"`
	CreatedAt *time.Time  `json:"created_at,omitempty"`
}
