package queue_test

import (
	"context"
	"testing"

	"github.com/felipefbs/testcontainers/connections/queue"
	"github.com/stretchr/testify/assert"
	testcontainers "github.com/testcontainers/testcontainers-go"
	kafkaContainer "github.com/testcontainers/testcontainers-go/modules/kafka"
)

func Test_Kafka(t *testing.T) {
	ctx := context.Background()

	kafkaContainer, err := kafkaContainer.RunContainer(
		ctx,
		kafkaContainer.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
	)
	if err != nil {
		panic(err)
	}

	bootstrapServers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		panic(err)
	}

	receiveChan := make(chan []byte, 10)
	messageReader := queue.NewMessageReader(bootstrapServers[0], map[string]string{"auto.offset.reset": "earliest"})
	go messageReader.ReadMessage(ctx, "testcontainers", receiveChan)
	assert.Nil(t, err)

	messageSender := queue.NewMessageSender(ctx, bootstrapServers[0], map[string]string{})
	err = messageSender.SendMessage(ctx, "testcontainers", []byte("message"), nil)
	assert.Nil(t, err)

	receivedMessage := <-receiveChan
	assert.Equal(t, []byte("message"), receivedMessage)

	// Clean up the container after
	defer func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()
}
