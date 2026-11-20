package protocol

import (
	"bytes"
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
