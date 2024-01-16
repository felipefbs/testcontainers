package queue

import (
	"context"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageSender struct {
	producer *kafka.Writer
	kafka    *kafka.Conn
}

func NewMessageSender(ctx context.Context, url string, additionalConfig map[string]string) *MessageSender {
	conn, err := kafka.DialLeader(ctx, "tcp", url, "testcontainers", 0)
	if err != nil {
		slog.Error("failed to connect to kafka", "error", err)
	}
	return &MessageSender{kafka: conn}
}

func (queue *MessageSender) SendMessage(ctx context.Context, topic string, message, key []byte) error {
	queue.kafka.SetWriteDeadline(time.Now().Add(10 * time.Second))
	msg := kafka.Message{
		Topic:     topic,
		Partition: 0,
		Key:       nil,
		Value:     message,
	}

	i, err := queue.kafka.WriteMessages(msg)
	if err != nil {
		slog.Error("failed to write message", "error", err)

		return err
	}

	slog.Info("bytes writen", "int", i)
	return nil
}

type MessageReader struct {
	consumer *kafka.Reader
}

func NewMessageReader(url string, additionalConfig map[string]string) *MessageReader {
	config := kafka.ReaderConfig{
		Brokers:     []string{url},
		Topic:       "testcontainers",
		GroupID:     "testcontainers-reader",
		StartOffset: kafka.FirstOffset,
		MaxBytes:    10e6,
	}

	r := kafka.NewReader(config)

	return &MessageReader{consumer: r}
}

func (queue *MessageReader) ReadMessage(ctx context.Context, topic string, receiveChan chan []byte) error {
	for {
		msg, err := queue.consumer.ReadMessage(ctx)
		if err != nil {
			slog.Error("failed to receive message", "error", err)
		}
		slog.Info("message received", "message", msg)
		receiveChan <- msg.Value
	}
}
