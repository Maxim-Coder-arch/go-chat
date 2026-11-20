package protocol

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	input := `{"type":"message","text":"Hello!"}` + "\n"

	decoder := NewDecoder(strings.NewReader(input))

	message, err := decoder.Decode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if message.Type != MessageTypeMessage {
		t.Errorf("expected type %q, got %q",
			MessageTypeMessage, message.Type)
	}

	if message.Text != "Hello!" {
		t.Errorf("expected text %q, got %q",
			"Hello!", message.Text)
	}
}

func TestDecoder_InvalidJSON(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("{invalid json}\n"))

	_, err := decoder.Decode()
	if err == nil {
		t.Fatal("expected decoding error, got nil")
	}
}

func TestDecoder_EOF(t *testing.T) {
	decoder := NewDecoder(strings.NewReader(""))

	_, err := decoder.Decode()

	if !errors.Is(err, io.EOF) {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestEncoder_Encode(t *testing.T) {
	var buffer bytes.Buffer

	encoder := NewEncoder(&buffer)

	message := ServerMessage{
		Type:     MessageTypeMessage,
		Username: "Maxim",
		Text:     "Hello!",
	}

	if err := encoder.Encode(message); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasSuffix(buffer.String(), "\n") {
		t.Error("expected newline at the end of encoded message")
	}
}

func TestClientMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message ClientMessage
		wantErr error
	}{
		{
			name: "valid message",
			message: ClientMessage{
				Type: MessageTypeMessage,
				Text: "Hello!",
			},
		},
		{
			name: "valid join",
			message: ClientMessage{
				Type: MessageTypeJoin,
			},
		},
		{
			name: "server-only type",
			message: ClientMessage{
				Type: MessageTypeSystem,
			},
			wantErr: ErrInvalidMessageType,
		},
		{
			name: "unknown type",
			message: ClientMessage{
				Type: "unknown",
			},
			wantErr: ErrInvalidMessageType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v",
					tt.wantErr, err)
			}
		})
	}
}


func TestClientEncoder_Encode(t *testing.T) {
	var buffer bytes.Buffer

	encoder := NewClientEncoder(&buffer)

	message := ClientMessage{
		Type:     MessageTypeJoin,
		Username: "Maxim",
	}

	if err := encoder.Encode(message); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	if !strings.HasSuffix(buffer.String(), "\n") {
		t.Fatal("expected newline at the end of encoded message")
	}

	var decoded ClientMessage

	if err := json.Unmarshal(
		bytes.TrimSpace(buffer.Bytes()),
		&decoded,
	); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.Type != MessageTypeJoin {
		t.Errorf("Type = %q, want %q", decoded.Type, MessageTypeJoin)
	}

	if decoded.Username != "Maxim" {
		t.Errorf("Username = %q, want Maxim", decoded.Username)
	}
}

func TestServerDecoder_Decode(t *testing.T) {
	input := `{"type":"system","text":"Welcome to the chat"}` + "\n"

	decoder := NewServerDecoder(strings.NewReader(input))

	message, err := decoder.Decode()
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if message.Type != MessageTypeSystem {
		t.Errorf("Type = %q, want %q", message.Type, MessageTypeSystem)
	}

	if message.Text != "Welcome to the chat" {
		t.Errorf("Text = %q, want Welcome to the chat", message.Text)
	}
}

func TestServerDecoder_InvalidJSON(t *testing.T) {
	decoder := NewServerDecoder(strings.NewReader("{invalid json}\n"))

	_, err := decoder.Decode()

	if err == nil {
		t.Fatal("expected JSON decoding error")
	}
}

func TestServerDecoder_EOF(t *testing.T) {
	decoder := NewServerDecoder(strings.NewReader(""))

	_, err := decoder.Decode()

	if !errors.Is(err, io.EOF) {
		t.Errorf("Decode() error = %v, want io.EOF", err)
	}
}