package protocol

import "time"

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

type ServerMessage struct {
    Type      MessageType `json:"type"`
    ID        string      `json:"id,omitempty"`
    Username  string      `json:"username,omitempty"`
    Text      string      `json:"text,omitempty"`
    CreatedAt *time.Time  `json:"created_at,omitempty"`
}