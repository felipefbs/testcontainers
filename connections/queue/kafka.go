package queue

import (
	"context"
	"log/slog"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MessageSender struct {
	producer *kafka.Producer
}

func NewMessageSender(url ...string) *MessageSender {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(url, ""),
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		slog.Error("failed to create producer", err)

		return nil
	}

	return &MessageSender{producer: p}
}

func (queue *MessageSender) SendMessage(ctx context.Context, topic string, message, key []byte, deliveryChan chan kafka.Event) error {
	msg := &kafka.Message{
		Value: message,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key: key,
	}

	err := queue.producer.Produce(msg, deliveryChan)
	if err != nil {
		slog.Error("Failed to produce message", err)

		return err
	}

	return nil
}

type MessageReader struct {
	consumer *kafka.Consumer
}

func NewMessageReader(url string, additionalConfig map[string]string) *MessageReader {
	config := &kafka.ConfigMap{
		"bootstrap.servers": url,
		"group.id":          "testcontainers-consumer",
	}

	for k, v := range additionalConfig {
		config.SetKey(k, v)
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		slog.Error("failed to create consumer")

		return nil
	}

	return &MessageReader{consumer: consumer}
}

func (queue *MessageReader) ReadMessage(ctx context.Context, topic string, receiveChan chan []byte) error {
	err := queue.consumer.Subscribe(topic, nil)
	if err != nil {
		slog.Error("failed to subscribe", "topic", topic)
		return err
	}

	go func() {
		for {
			ev := queue.consumer.Poll(1000)
			switch e := ev.(type) {
			case *kafka.Message:
				slog.Info("message received", "value", e.Value)
				receiveChan <- e.Value
			case kafka.Error:
				slog.Error("failed to receive message")
			}

		}
	}()

	return nil
}
