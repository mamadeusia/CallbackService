package nats

import (
	"context"
	"time"

	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

type PushBasedTopic struct {
	pushBasedConsumer events.Stream

	topic   string
	handler EventHandler
	autoAck bool

	//if we need to create delayed and topics with delayed we need to set this
	readWithDelay bool
	consumeDelay  time.Duration

	// we can receive events immediately after the publish but if we wnat to batch events we can use this batch
	// but it is preferred for the situations where you have only one instance of the service .
	// if you have more than one consumer services for avoid ack issues use pullbased consumer instead of push based consumer.
	// this batch vs pull based consumes more memory and cpu but it is you can receive your batch events faster.
	batchEnabled bool
	batchSize    int
}

type PushBasedTopicConfiguration func(p *PushBasedTopic) error

func (pbt *PushBasedTopic) StartConsume(ctx context.Context) error {
	eventChan, err := pbt.pushBasedConsumer.Consume(pbt.topic, consumeOptionsFromPushBasedTopic(pbt)...)
	_ = eventChan
	if err != nil {
		logger.Error(err)
		return err
	}
	go func() {
		for keepGoing := true; keepGoing; {
			for {
				select {
				case <-ctx.Done():
					keepGoing = false
				case event, ok := <-eventChan:
					if !ok {
						logger.Errorf("consumer channel closed, topic :: %s ", pbt.topic)
						keepGoing = false
					}
					if pbt.batchEnabled {

					} else {
						if pbt.readWithDelay {
							currentTime := time.Now()
							if currentTime.After(event.Timestamp.Add(pbt.consumeDelay)) {
								sliceEvent := []*events.Event{&event}
								if err := pbt.handler(ctx, sliceEvent); err != nil {
									logger.Error(err)
								}
							}

						} else {
							sliceEvent := []*events.Event{&event}
							if err := pbt.handler(ctx, sliceEvent); err != nil {
								logger.Error(err)
							}
						}
					}
				}
			}
		}
	}()
	return nil
}

func consumeOptionsFromPushBasedTopic(pbt *PushBasedTopic) []events.ConsumeOption {
	return nil
}
