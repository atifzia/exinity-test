package kafka

import (
	"context"
	"fmt"
	"log"
	"payment-gateway/internal/util"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sony/gobreaker"
)

// CircuitBreaker configuration for Kafka publisher
var cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
	Name:        "KafkaPublisher",
	MaxRequests: 1,               // Allows 1 request at a time
	Interval:    5 * time.Second, // Reset the breaker every 5 seconds
	Timeout:     3 * time.Second, // Timeout for an operation
})

type (
	IProducer interface {
		PubTx(ctx context.Context, txId int64, msg []byte, format string) error
		Close() error
	}

	Producer struct {
		writer *kafka.Writer
	}
)

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(kafkaURL string) IProducer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(kafkaURL),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
			BatchTimeout:           10 * time.Millisecond,
		},
	}
}

// PubTx publishes a message to the Kafka topic
func (p Producer) PubTx(ctx context.Context, txId int64, msg []byte, format string) error {
	if p.writer == nil {
		return fmt.Errorf("kafka writer is not initialized")
	}

	// Circuit breaker logic for Kafka publishing
	return util.RetryOperation(func() error {
		return PubTxWithCB(func() error {
			topic, err := GetTopic(format)
			if err != nil {
				return err
			}

			kafkaMessage := kafka.Message{
				Key:   []byte(fmt.Sprintf("%d", txId)),
				Value: msg,
				Topic: topic,
			}

			err = p.writer.WriteMessages(ctx, kafkaMessage)
			if err != nil {
				log.Printf("failed to publish msg to Kafka: %v", err)
				return err
			}

			log.Println("message successfully published to Kafka on topic " + string(topic))

			return nil
		})
	}, 5)
}

// PubTxWithCB uses a circuit breaker to manage Kafka publishing
func PubTxWithCB(op func() error) error {
	_, err := cb.Execute(func() (interface{}, error) {
		return nil, op()
	})
	return err
}

// GetTopic returns the appropriate Kafka topic based on the data format
func GetTopic(dataFormat string) (string, error) {
	switch dataFormat {
	case "application/json":
		return "transactions.json", nil
	case "text/xml", "application/xml":
		return "transactions.soap", nil
	default:
		return "", fmt.Errorf("unsupported data format: %s", dataFormat)
	}
}

// Close the writer when the system shuts down
func (p Producer) Close() error {
	return p.writer.Close()
}
