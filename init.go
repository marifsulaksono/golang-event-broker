package eventbroker

import (
	"errors"

	"github.com/marifsulaksono/golang-event-broker/config"
	"github.com/marifsulaksono/golang-event-broker/constant"
	"github.com/marifsulaksono/golang-event-broker/event"
	"github.com/marifsulaksono/golang-event-broker/kafka"
	"github.com/marifsulaksono/golang-event-broker/rabbitmq"
)

// New initializes and returns an implementation of the EventBus interface
// based on the configured message broker.
//
// cfg is Configuration struct containing broker type and specific settings for each broker.
//
// Note:
//   - The actual broker to be used is determined by cfg.Broker.
//   - Only the required broker's implementation is initialized.
func New(cfg config.Config) (event.EventBus, error) {
	switch cfg.Broker {
	case constant.RabbitMQBroker:
		return rabbitmq.New(cfg.RabbitMQ)
	case constant.KafkaBroker:
		return kafka.New(cfg.Kafka)
	default:
		return nil, errors.New("unsupported broker: " + cfg.Broker)
	}
}
