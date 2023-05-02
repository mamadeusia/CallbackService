package nats

import (
	"go-micro.dev/v4/events"
)

func WithPullBasedStore(es events.Store) NatsConfiguration {
	return func(h *NatsJS) error {
		h.pullBasedConsumer = es
		return nil
	}
}

func WithPushBasedStream(es events.Stream) NatsConfiguration {
	return func(h *NatsJS) error {
		h.pushBasedConsumer = es
		return nil
	}
}

func WithPullBasedTopic(pbt *PullBasedTopic) NatsConfiguration {
	return func(h *NatsJS) error {
		if h.topics == nil {
			h.topics = []Topic{}
		}
		if h.pullBasedConsumer == nil {
			return ErrPullBasedConsumerNotFound
		}
		//here we need to set the consumer of the topic
		pbt.pullBasedConsumer = h.pullBasedConsumer
		h.topics = append(h.topics, pbt)
		return nil
	}
}

func WithPushBasedTopic(pbt *PushBasedTopic) NatsConfiguration {
	return func(h *NatsJS) error {
		if h.topics == nil {
			h.topics = []Topic{}
		}
		if h.pushBasedConsumer == nil {
			return ErrPushBasedConsumerNotFound
		}
		//here we need to set the consumer of the topic
		pbt.pushBasedConsumer = h.pushBasedConsumer
		h.topics = append(h.topics, pbt)
		return nil
	}
}
