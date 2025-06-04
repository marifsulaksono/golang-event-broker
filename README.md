# Go Event Broker Library
A flexible event-driven messaging library for Go microservices, supporting RabbitMQ and Kafka brokers with a pluggable architecture. Designed for easy switching between brokers by changing configuration only, no code changes required.

## Features
- Unified interface for publishing and subscribing to events
- Support for RabbitMQ and Kafka brokers
- Simple config-based broker selection
- JSON payload encoding/decoding
- Context-aware publish and subscribe methods
- Automatic message consumption with handler callback
- Graceful shutdown on context cancellation (Kafka)

## Installation
```bash
go get github.com/marifsulaksono/golang-event-broker
```

## Usage
### Configuration
Define your config structs, for example:
```go
import (
    "os"

    ebconfig "github.com/marifsulaksono/golang-event-broker/config"
    ebconst  "github.com/marifsulaksono/golang-event-broker/constant"
)

func loadConfig() ebconfig.Config {
	broker := os.Getenv("BROKER")

	cfg := ebconfig.Config{
		Broker: broker,
	}

	if broker == ebconst.RabbitMQBroker {
		cfg.RabbitMQ = &ebconfig.RabbitMQConfig{
			URI:          os.Getenv("RABBITMQ_URI"),
			Exchange:     os.Getenv("RABBITMQ_EXCHANGE"),
			ExchangeType: os.Getenv("RABBITMQ_EXCHANGE_TYPE"),
            Durable:      true,
		}
	}

	if broker == ebconst.KafkaBroker {
		cfg.Kafka = &ebconfig.KafkaConfig{
			Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
			GroupID: os.Getenv("KAFKA_GROUP_ID"),
		}
	}

	return cfg
}
```

### Initialization
Create an event bus instance by passing your config:
```go
import (
    "log"

    eventbroker "github.com/marifsulaksono/golang-event-broker"
)

eb, err := eventbroker.New(config)
if err != nil {
    log.Fatal(err)
}
```

### Publishing Messages
```go
type EmailConfig struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

email := EmailConfig{
    To:      fmt.Sprintf("user%d@example.com", i),
    Subject: fmt.Sprintf("Notification %d", i),
    Body:    fmt.Sprintf("This is notification email number %d sent at %s", i, time.Now().Format(time.RFC3339)),
}

err = eb.Publish(ctx, "notification.email", email)
if err != nil {
    log.Println("Publish failed:", err)
} else {
    log.Println("Published email:", email)
}
```

### Subscribing to Messages
```go
type EmailHandler struct{}

func (h *EmailHandler) Handle(ctx context.Context, msg any) error {
    data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var email EmailConfig
	if err := json.Unmarshal(data, &email); err != nil {
		return err
	}

	log.Printf("Sending email to %s with subject %s\n", email.To, email.Subject)
	return nil
}

handler := &EmailHandler{}
err = eb.Subscribe(ctx, "notification.email", handler)
if err != nil {
    log.Fatal("Failed subscribe:", err)
}

// Wait until termination signal
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
<-c
log.Println("Shutting down...")
```

## How It Works
- Publish: Encodes payload as JSON and sends it to the broker exchange or topic.
- Subscribe: Starts a background goroutine consuming messages from broker queue/topic and calls the handler.
- Broker Agnostic: You can switch between RabbitMQ and Kafka by changing Broker config without code changes.

## License
MIT License