package nats

import (
	"context"
	"errors"

	"go-micro.dev/v4/events"
)

var (
	ErrPullBasedConsumerNotFound = errors.New("nats : pull-based consumer not found")
	ErrPushBasedConsumerNotFound = errors.New("nats : push-based consumer not found")
)

type EventHandler func(ctx context.Context, e []*events.Event) error

type Topic interface {
	StartConsume(context.Context) error
}

type NatsJS struct {
	pullBasedConsumer events.Store
	pushBasedConsumer events.Stream
	topics            []Topic
}

type NatsConfiguration func(a *NatsJS) error

func New(cfgs ...NatsConfiguration) (*NatsJS, error) {
	natsjs := &NatsJS{}
	for _, cfg := range cfgs {
		err := cfg(natsjs)
		if err != nil {
			return nil, err
		}
	}
	return natsjs, nil
}

func (n *NatsJS) Start(ctx context.Context) error {

	for _, topic := range n.topics {
		if err := topic.StartConsume(ctx); err != nil {
			return err
		}
	}

	return nil
}
