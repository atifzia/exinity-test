package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type MockKafkaProducer struct {
	Messages []kafka.Message // Captures the sent messages
}

// PubTx simulates sending messages to Kafka
func (p *MockKafkaProducer) PubTx(ctx context.Context, txId int64, msg []byte, format string) error {
	log.Printf("Mock Kafka Producer sent message: %s\n", string(msg))
	p.Messages = append(p.Messages, kafka.Message{})
	return nil
}

// Close simulates closing the Kafka producer
func (p *MockKafkaProducer) Close() error {
	log.Println("Mock Kafka Producer closed")
	return nil
}
