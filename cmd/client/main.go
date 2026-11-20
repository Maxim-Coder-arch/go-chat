package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"chat-app/internal/protocol"
	"chat-app/internal/transport/cli"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter username: ")

	if !reader.Scan() {
		log.Fatal("failed to read username")
	}

	username := strings.TrimSpace(reader.Text())

	if username == "" {
		log.Fatal("username cannot be empty")
	}

	client, err := cli.NewClient("localhost:8080")
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer client.Close()

	if err := client.Send(protocol.ClientMessage{
		Type:     protocol.MessageTypeJoin,
		Username: username,
	}); err != nil {
		log.Fatalf("failed to join chat: %v", err)
	}

	fmt.Println("Connected to chat!")
	fmt.Println("Type a message or /quit to exit.")
	fmt.Println("--------------------------------")

	// Receive messages from the server.
	go func() {
		for {
			message, err := client.Receive()
			if err != nil {
				fmt.Printf("\nConnection closed: %v\n", err)
				return
			}

			switch message.Type {
			case protocol.MessageTypeMessage:
				fmt.Printf("\n[%s]: %s\n", message.Username, message.Text)

			case protocol.MessageTypeSystem:
				fmt.Printf("\n[SYSTEM]: %s\n", message.Text)

			case protocol.MessageTypeError:
				fmt.Printf("\n[ERROR]: %s\n", message.Text)
			}

			fmt.Print("> ")
		}
	}()

	// Send messages from terminal.
	for {
		fmt.Print("> ")

		if !reader.Scan() {
			break
		}

		text := strings.TrimSpace(reader.Text())

		if text == "" {
			continue
		}

		if text == "/quit" {
			break
		}

		if err := client.Send(protocol.ClientMessage{
			Type: protocol.MessageTypeMessage,
			Text: text,
		}); err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			break
		}
	}

	_ = client.Send(protocol.ClientMessage{
		Type: protocol.MessageTypeLeave,
	})

	fmt.Println("Disconnected from chat.")

	if err := reader.Err(); err != nil {
		log.Printf("input error: %v", err)
	}
}