package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/felipefbs/testcontainers/connections/database"
	"github.com/felipefbs/testcontainers/connections/queue"
	"github.com/felipefbs/testcontainers/message"
)

func main() {
	db := database.NewConnection(os.Getenv("DATABASE_URL"))
	KafkaConsume(context.Background(), db)
}

func KafkaConsume(ctx context.Context, db *database.Database) {
	reader := queue.NewMessageReader(os.Getenv("KAFKA_SERVER_URL"), map[string]string{})
	receiveChan := make(chan []byte, 10)
	err := reader.ReadMessage(ctx, "testcontainers", receiveChan)
	if err != nil {
		slog.Error("failed to read message", "error", err)
	}

	for {
		receivedMsg := <-receiveChan
		msg := message.Message{}
		err := json.Unmarshal(receivedMsg, &msg)
		if err != nil {
			slog.Error("failed to decode message", "error", err)
		}
		msg.ReceivedAt = time.Now()
		err = db.Save(ctx, msg)
		if err != nil {
			slog.Error("failed to save message", "error", err)
		}
	}
}
