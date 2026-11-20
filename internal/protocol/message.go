package protocol

import (
    "errors"
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
    Type MessageType `json:"type"`
    Text string      `json:"text,omitempty"`
}

func (m ClientMessage) Validate() error {
    switch m.Type {
    case MessageTypeJoin, MessageTypeMessage, MessageTypeLeave:
        return nil
    default:
        return ErrInvalidMessageType
    }
}

type ServerMessage struct {
    Type      MessageType `json:"type"`
    ID        string      `json:"id,omitempty"`
    Username  string      `json:"username,omitempty"`
    Text      string      `json:"text,omitempty"`
    CreatedAt *time.Time  `json:"created_at,omitempty"`
}