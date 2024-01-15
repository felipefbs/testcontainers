package message

import "time"

type Message struct {
	ID         int       `json:"id"`
	SendedAt   time.Time `json:"sendedAt"`
	ReceivedAt time.Time `json:"receivedAt"`
	InsertedAt time.Time `json:"insertedAt"`
	Message    string    `json:"message"`
}
