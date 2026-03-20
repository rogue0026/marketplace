package broker

import (
	"context"
	"encoding/json"
	"errors"
	"order_service/internal/service"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	ErrEventWriterNotFound = errors.New("event writer not found")
)

type KafkaEventPublisher struct {
	writers map[string]*kafka.Writer
}

func NewKafkaEventPublisher(brokers []string, events []service.Event) *KafkaEventPublisher {
	writers := make(map[string]*kafka.Writer)
	for _, e := range events {
		eventName := e.GetEventName()
		writer := &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  eventName,
			RequiredAcks:           kafka.RequireAll,
			Async:                  true,
			Completion:             nil,
			Compression:            kafka.Snappy,
			AllowAutoTopicCreation: true, // todo
		}
		writers[eventName] = writer
	}

	ep := &KafkaEventPublisher{
		writers: writers,
	}

	return ep
}

func (ep *KafkaEventPublisher) Publish(ctx context.Context, e service.Event) error {
	writer, ok := ep.writers[e.GetEventName()]
	if !ok {
		return ErrEventWriterNotFound
	}

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: e.GetEventName(),
		Value: data,
		Time:  time.Now(),
	}

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
