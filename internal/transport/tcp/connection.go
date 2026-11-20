package tcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"

	"chat-app/internal/application/chat"
	"chat-app/internal/domain"
	"chat-app/internal/protocol"
)

const (
	ClientSendBuffer = 32
)

type Connection struct {
	conn    net.Conn
	decoder *protocol.Decoder
	encoder *protocol.Encoder

	send chan protocol.ServerMessage

	username string

	handler chat.MessageHandler
	hub     chat.ClientRegistry

	closeOnce sync.Once

	done chan struct{}
}

func NewConnection(
	conn net.Conn,
	handler chat.MessageHandler,
	hub chat.ClientRegistry,
) *Connection {
	return &Connection{
		conn:    conn,
		decoder: protocol.NewDecoder(conn),
		encoder: protocol.NewEncoder(conn),
		send:    make(chan protocol.ServerMessage, ClientSendBuffer),
		done:    make(chan struct{}),
		handler: handler,
		hub:     hub,
	}
}

var ErrClientQueueFull = errors.New("client send queue is full")

func (c *Connection) Send(message domain.Message) error {
	serverMessage := protocol.ServerMessage{
		Type:      protocol.MessageTypeMessage,
		ID:        message.ID,
		Username:  message.Username,
		Text:      message.Text,
		CreatedAt: &message.CreatedAt,
	}

	select {
	case <-c.done:
		return net.ErrClosed
	default:
	}

	select {
	case c.send <- serverMessage:
		return nil
	case <-c.done:
		return net.ErrClosed
	default:
		return ErrClientQueueFull
	}
}

func (c *Connection) Close() error {
	var err error

	c.closeOnce.Do(func() {
		close(c.done)
		err = c.conn.Close()
	})

	return err
}

func (c *Connection) readLoop(ctx context.Context) error {
	joinMessage, err := c.decoder.Decode()
	if err != nil {
		return fmt.Errorf("read join message: %w", err)
	}

	if err := joinMessage.Validate(); err != nil {
		return err
	}

	if joinMessage.Type != protocol.MessageTypeJoin {
		return errors.New("first message must be join")
	}

	c.username = strings.TrimSpace(joinMessage.Username)

	if c.username == "" {
		return domain.ErrEmptyUsername
	}

	c.hub.Register(c)

	for {
		message, err := c.decoder.Decode()
		if err != nil {
			return err
		}

		if err := message.Validate(); err != nil {
			return err
		}

		switch message.Type {
		case protocol.MessageTypeMessage:
			_, err := c.handler.SendMessage(
				ctx,
				c.username,
				message.Text,
			)

			if err != nil {
				if sendErr := c.sendError(err.Error()); sendErr != nil {
					return sendErr
				}
			}

		case protocol.MessageTypeLeave:
			return nil

		default:
			return errors.New("unsupported message type")
		}
	}
}

func (c *Connection) writeLoop() error {
	for {
		select {
		case <-c.done:
			return nil

		case message := <-c.send:
			if err := c.encoder.Encode(message); err != nil {
				return fmt.Errorf("write message: %w", err)
			}
		}
	}
}

func (c *Connection) sendError(text string) error {
	select {
	case <-c.done:
		return net.ErrClosed
	case c.send <- protocol.ServerMessage{
		Type: protocol.MessageTypeError,
		Text: text,
	}:
		return nil
	default:
		return ErrClientQueueFull
	}
}

func (c *Connection) Run(ctx context.Context) {
	defer c.Close()

	stop := context.AfterFunc(ctx, func() {
		_ = c.Close()
	})
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := c.writeLoop(); err != nil {
			_ = c.Close()
		}
	}()

	if err := c.readLoop(ctx); err != nil &&
		!errors.Is(err, io.EOF) &&
		!errors.Is(err, net.ErrClosed) {
		slog.Warn("Client connection ended", "error", err)
	}

	c.hub.Unregister(c)

	wg.Wait()
}
