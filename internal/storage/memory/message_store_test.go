package memory

import (
	"context"
	"testing"
	"time"

	"chat-app/internal/domain"
)

func createMessage(id string) domain.Message {
	return domain.Message{
		ID:        id,
		Username:  "tester",
		Text:      "Message " + id,
		CreatedAt: time.Now().UTC(),
	}
}

func TestMessageStore_SaveAndGetRecent(t *testing.T) {
	store := NewMessageStore(3)
	ctx := context.Background()

	for i := 1; i <= 4; i++ {
		message := createMessage(string(rune('0' + i)))

		if err := store.Save(ctx, message); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	messages, err := store.GetRecent(ctx, 3)
	if err != nil {
		t.Fatalf("GetRecent() error = %v", err)
	}

	if len(messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(messages))
	}

	expectedIDs := []string{"2", "3", "4"}

	for i, message := range messages {
		if message.ID != expectedIDs[i] {
			t.Errorf(
				"message[%d].ID = %q, want %q",
				i,
				message.ID,
				expectedIDs[i],
			)
		}
	}
}

func TestMessageStore_GetRecentLimit(t *testing.T) {
	store := NewMessageStore(10)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		if err := store.Save(ctx, createMessage(string(rune('0'+i)))); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	messages, err := store.GetRecent(ctx, 2)
	if err != nil {
		t.Fatalf("GetRecent() error = %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	if messages[0].ID != "4" || messages[1].ID != "5" {
		t.Errorf("unexpected messages: %+v", messages)
	}
}

func TestMessageStore_GetRecentReturnsCopy(t *testing.T) {
	store := NewMessageStore(10)
	ctx := context.Background()

	if err := store.Save(ctx, createMessage("1")); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	messages, err := store.GetRecent(ctx, 10)
	if err != nil {
		t.Fatalf("GetRecent() error = %v", err)
	}

	messages[0].Text = "Modified"

	stored, err := store.GetRecent(ctx, 10)
	if err != nil {
		t.Fatalf("GetRecent() error = %v", err)
	}

	if stored[0].Text == "Modified" {
		t.Fatal("modifying returned messages affected internal storage")
	}
}

func TestMessageStore_SaveWithCancelledContext(t *testing.T) {
	store := NewMessageStore(10)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := store.Save(ctx, createMessage("1"))

	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestMessageStore_GetRecentWithCancelledContext(t *testing.T) {
	store := NewMessageStore(10)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.GetRecent(ctx, 10)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestMessageStore_DefaultLimit(t *testing.T) {
	store := NewMessageStore(0)

	if store.limit != DefaultHistoryLimit {
		t.Errorf(
			"limit = %d, want %d",
			store.limit,
			DefaultHistoryLimit,
		)
	}
}