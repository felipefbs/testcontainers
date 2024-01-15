package main

import (
	"context"
	"encoding/json"
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
	reader.ReadMessage(ctx, "testcontainers", receiveChan)

	for {
		receivedMsg := <-receiveChan
		msg := message.Message{}
		json.Unmarshal(receivedMsg, &msg)
		msg.ReceivedAt = time.Now()
		db.Save(ctx, msg)
	}
}
