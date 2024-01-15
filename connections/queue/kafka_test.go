package queue_test

import (
	"context"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	receiveChan := make(chan []byte)
	messageReader := queue.NewMessageReader(bootstrapServers[0], map[string]string{"auto.offset.reset": "earliest"})
	err = messageReader.ReadMessage(ctx, "testcontainers", receiveChan)
	assert.Nil(t, err)

	deliveryChan := make(chan kafka.Event)
	messageSender := queue.NewMessageSender(bootstrapServers...)
	err = messageSender.SendMessage(ctx, "testcontainers", []byte("message"), nil, deliveryChan)
	assert.Nil(t, err)

	e := <-deliveryChan
	m := e.(*kafka.Message)
	assert.NotNil(t, m.TopicPartition)
	t.Log(m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	close(deliveryChan)

	receivedMessage := <-receiveChan
	assert.Equal(t, []byte("message"), receivedMessage)

	// Clean up the container after
	defer func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()
}
