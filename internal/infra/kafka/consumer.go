package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kafkaGo "github.com/segmentio/kafka-go"
)

type Consumer interface {
	Consume(ctx context.Context) error
	Close() error
}

type kafkaConsumer struct {
	reader  *kafkaGo.Reader
	handler func(context.Context, Event) error
}

func NewConsumer(cfg ConsumerConfig, handler func(context.Context, Event) error) (Consumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers are required")
	}
	if cfg.Topic == "" {
		return nil, fmt.Errorf("kafka topic is required")
	}
	if cfg.GroupID == "" {
		return nil, fmt.Errorf("kafka group id is required")
	}

	r := kafkaGo.NewReader(kafkaGo.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		Dialer:         &kafkaGo.Dialer{Timeout: cfg.DialTimeout},
	})

	return &kafkaConsumer{reader: r, handler: handler}, nil
}

func (c *kafkaConsumer) Consume(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var event Event
		if err := json.Unmarshal(m.Value, &event); err != nil {
			return fmt.Errorf("unmarshal kafka event: %w", err)
		}

		if err := c.handler(ctx, event); err != nil {
			return err
		}
	}
}

func (c *kafkaConsumer) Close() error {
	if c.reader == nil {
		return nil
	}
	return c.reader.Close()
}
