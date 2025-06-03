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
type RabbitMQConfig struct {
    URI          string
    Exchange     string
    ExchangeType string
    QueueName    string
    RoutingKey   string
    Durable      bool
    AutoDelete   bool
}

type KafkaConfig struct {
    Brokers []string
    GroupID string
}

type Config struct {
    Broker   string // "rabbitmq" or "kafka"
    RabbitMQ RabbitMQConfig
    Kafka    KafkaConfig
}
```

### Initialization
Create an event bus instance by passing your config:
```go
import "github.com/yourusername/go-event-bus"

bus, err := eventbroker.New(config)
if err != nil {
    log.Fatal(err)
}
```

### Publishing Messages
```go
err := bus.Publish(ctx, "topic.name", map[string]interface{}{
    "user_id": 123,
    "action":  "created",
})
if err != nil {
    log.Println("Publish error:", err)
}
```

### Subscribing to Messages
```go
type MyHandler struct{}

func (h *MyHandler) Handle(ctx context.Context, event map[string]interface{}) error {
    fmt.Println("Received event:", event)
    return nil
}

handler := &MyHandler{}
err = bus.Subscribe(ctx, "topic.name", handler)
if err != nil {
    log.Fatal(err)
}
```

## How It Works
- Publish: Encodes payload as JSON and sends it to the broker exchange or topic.
- Subscribe: Starts a background goroutine consuming messages from broker queue/topic and calls the handler.
- Broker Agnostic: You can switch between RabbitMQ and Kafka by changing Broker config without code changes.

## License
MIT License