package domain

import "time"

type Message struct {
	ID string
	Username string
	Text string
	CreatedAt time.Time
}