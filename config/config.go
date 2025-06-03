package config

type Config struct {
	Broker   string          // name of the message broker, ex "rabbitmq"
	RabbitMQ *RabbitMQConfig // configuration for RabbitMQ
	Kafka    *KafkaConfig    // configuration for Kafka
}

type RabbitMQConfig struct {
	URI          string // RabbitMQ connection URI to use for connecting
	Exchange     string // Exchange name for publish/subscribe
	ExchangeType string // Exchange type to use for routing message delivery, ex ("direct", "topic", "fanout", "headers")

	// Optional
	QueueName  string // Name of the queue to consume from (used by subscriber).
	RoutingKey string // Routing key to publish with (used by publisher).
	Durable    bool   // Whether the exchange/queue should survive broker restarts.
	AutoDelete bool   // Whether the exchange/queue should be deleted when the last consumer unsubscribes.
}

type KafkaConfig struct {
	Brokers []string // List of Kafka brokers to connect to
	GroupID string   // Kafka consumer group ID
}
