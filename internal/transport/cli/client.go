package cli

import (
	"fmt"
	"net"

	"chat-app/internal/protocol"
)

type Client struct {
	conn    net.Conn
	encoder *protocol.ClientEncoder
	decoder *protocol.ServerDecoder
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("connect to server: %w", err)
	}

	return &Client{
		conn:    conn,
		encoder: protocol.NewClientEncoder(conn),
		decoder: protocol.NewServerDecoder(conn),
	}, nil
}

func (c *Client) Send(message protocol.ClientMessage) error {
	if err := message.Validate(); err != nil {
		return fmt.Errorf("validate message: %w", err)
	}

	if err := c.encoder.Encode(message); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (c *Client) Receive() (protocol.ServerMessage, error) {
	message, err := c.decoder.Decode()
	if err != nil {
		return protocol.ServerMessage{}, fmt.Errorf("receive message: %w", err)
	}

	return message, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}