package domain

import (
	"errors"
	"testing"
)

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
        name    string
        message Message
        wantErr error
    }{
        {
            name: "valid message",
            message: Message{
                Username: "Maxim",
                Text:     "Hello, everyone!",
            },
            wantErr: nil,
        },
        {
            name: "empty username",
            message: Message{
                Username: " ",
                Text:     "Hello!",
            },
            wantErr: ErrEmptyUsername,
        },
        {
            name: "empty text",
            message: Message{
                Username: "Maxim",
                Text:     "   ",
            },
            wantErr: ErrEmptyMessage,
        },
    }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
            err := tt.message.Validate()

            if !errors.Is(err, tt.wantErr) {
                t.Errorf("expected error %v, got %v", tt.wantErr, err)
            }
        })
	}
}