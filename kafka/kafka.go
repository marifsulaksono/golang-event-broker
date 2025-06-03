package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/marifsulaksono/golang-event-broker/config"
	"github.com/marifsulaksono/golang-event-broker/event"

	"github.com/segmentio/kafka-go"
)

type KafkaBus struct {
	cfg    *config.KafkaConfig
	writer *kafka.Writer
}

func New(cfg *config.KafkaConfig) (event.EventBus, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers list cannot be empty")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaBus{cfg, writer}, nil
}

// Publish sends a message to the given Kafka topic.
// The payload will be marshaled to JSON before publishing.
func (b *KafkaBus) Publish(ctx context.Context, topic string, payload any) error {
	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return b.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: value,
	})
}

// Subscribe consumes messages from the specified Kafka topic and passes them
// to the provided EventHandler. It runs in a separate goroutine and respects context cancellation.
func (b *KafkaBus) Subscribe(ctx context.Context, topic string, handler event.EventHandler) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: b.cfg.Brokers,
		GroupID: b.cfg.GroupID,
		Topic:   topic,
	})

	go b.consumeLoop(ctx, r, handler)
	return nil
}

// consumeLoop continuously reads messages from the Kafka reader and delegates handling.
func (b *KafkaBus) consumeLoop(ctx context.Context, r *kafka.Reader, handler event.EventHandler) {
	defer func() {
		if err := r.Close(); err != nil {
			log.Println("Kafka reader close error:", err)
		}
	}()

	for {
		if ctx.Err() != nil {
			return // exit when context is cancelled
		}

		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		b.processMessage(ctx, msg, handler)
	}
}

// processMessage unmarshals the Kafka message and calls the handler.
func (b *KafkaBus) processMessage(ctx context.Context, msg kafka.Message, handler event.EventHandler) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		log.Println("Failed to unmarshal Kafka message:", err)
		return
	}

	if err := handler.Handle(ctx, data); err != nil {
		log.Println("Error handling Kafka message:", err)
	}
}
