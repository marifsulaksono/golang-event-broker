package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/marifsulaksono/golang-event-broker/config"
	"github.com/marifsulaksono/golang-event-broker/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBus struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *config.RabbitMQConfig
}

// New initializes a RabbitMQ event bus connection and prepares the exchange and optional queue.
func New(cfg *config.RabbitMQConfig) (event.EventBus, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	durable := cfg.Durable
	autoDelete := cfg.AutoDelete

	// Declare exchange
	err = ch.ExchangeDeclare(
		cfg.Exchange,
		cfg.ExchangeType,
		durable,
		autoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Declare queue if queue name is not empty
	if cfg.QueueName != "" {
		q, err := ch.QueueDeclare(
			cfg.QueueName,
			durable,
			autoDelete,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}

		// Bind queue ke exchange
		err = ch.QueueBind(
			q.Name,
			cfg.RoutingKey,
			cfg.Exchange,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	return &RabbitMQBus{
		conn:    conn,
		channel: ch,
		cfg:     cfg,
	}, nil
}

// Publish sends a message payload to the specified topic (routing key) on the configured RabbitMQ exchange.
//   - ctx: Context to allow cancellation and timeout control for the publishing operation.
//   - topic: The routing key to use for routing the message within the exchange.
//   - payload: The message content to be sent; it will be marshaled into JSON.
func (b *RabbitMQBus) Publish(ctx context.Context, topic string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return b.channel.PublishWithContext(ctx, b.cfg.Exchange, topic, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// Subscribe starts consuming messages from the configured RabbitMQ exchange and routing key (topic).
//   - ctx: Context for controlling cancellation and deadlines.
//   - topic: The routing key used to bind the queue to the exchange.
//   - handler: The event handler called for each incoming message.
func (b *RabbitMQBus) Subscribe(ctx context.Context, topic string, handler event.EventHandler) error {
	queueName := b.cfg.QueueName

	if queueName == "" {
		// create queue anonymous for subscriber ephemeral
		q, err := b.channel.QueueDeclare(
			"",
			false,
			true,
			true,
			false,
			nil,
		)
		if err != nil {
			return err
		}
		queueName = q.Name
	}

	// bind queue to exchange with routing key
	err := b.channel.QueueBind(queueName, topic, b.cfg.Exchange, false, nil)
	if err != nil {
		return err
	}

	msgs, err := b.channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// goroutine for listening messages
	go func() {
		for msg := range msgs {
			var m map[string]interface{}
			if err := json.Unmarshal(msg.Body, &m); err != nil {
				log.Println("error unmarshalling message:", err)
				continue
			}
			if err := handler.Handle(ctx, m); err != nil {
				log.Println("error handling message:", err)
			}
		}
	}()

	return nil
}
