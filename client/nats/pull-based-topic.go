package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/google/uuid"
	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

type PullBasedTopic struct {
	pullBasedConsumer events.Store

	topic    string
	duration time.Duration
	maxItems int
	handler  EventHandler
	autoAck  bool
	group    string

	//if we need to create delayed and topics with delayed we need to set this
	readWithDelay bool
	consumeDelay  time.Duration

	hastilyRecieveData            bool
	minReqiredItmesRecieveHastily int
	hastilyRecieveDataNotif       chan struct{}
}

// 1pull 80 -> 80-80
// callback.tenMin

// "callback.now"
//
//	1 2 3 4 5
//
// now - 1 - 15 - 25 - 1 -
type PullBasedTopicConfiguration func(p *PullBasedTopic) error

func NewPullBasedTopic(cfgs ...PullBasedTopicConfiguration) (*PullBasedTopic, error) {
	pullBasedTopic := &PullBasedTopic{
		maxItems: 1,
		group:    uuid.New().String(),
	}

	for _, cfg := range cfgs {
		err := cfg(pullBasedTopic)
		if err != nil {
			return nil, err
		}
	}
	return pullBasedTopic, nil
}

// integration test
func (pbt *PullBasedTopic) StartConsume(ctx context.Context) error {
	go func() {
		for keepGoing := true; keepGoing; {
			expire := time.After(pbt.duration)
			select {
			case <-ctx.Done():
				keepGoing = false
			case <-expire:
				pulledEvents, err := pbt.read(ctx)
				if err != nil {
					logger.Error(err)
				}
				if len(pulledEvents) > 0 {
					err = pbt.process(ctx, pulledEvents)
					if err != nil {
						logger.Error(err)
					}
				}
			case <-pbt.hastilyRecieveDataNotif:
				pulledEvents, err := pbt.readHastily(ctx)
				if err != nil {
					logger.Error(err)
				}
				err = pbt.processHastily(ctx, pulledEvents)
				if err != nil {
					logger.Error(err)
				}

			}
		}
	}()
	return nil
}

func readOptionsFromPullBasedTopic(pbt *PullBasedTopic) []events.ReadOption {
	var output []events.ReadOption
	output = append(output, events.ReadLimit(uint(pbt.maxItems)))

	// output = append(output, events.ReadLimit(uint(pbt.maxItems)))

	if pbt.group != "" {
		output = append(output, natsjs.WithReadGroupName(pbt.group))
	}

	//if we need to read with delay then we have to ack message after time checking
	if !pbt.readWithDelay {
		output = append(output, natsjs.WithReadAutoAck(pbt.autoAck))
	} else {
		output = append(output, natsjs.WithReadAutoAck(false))
	}
	// natsjs.WithReadAutoAck(pbt.autoAck),
	return output
}
func PPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
	return string(s)
}

func readOptionsFromPullBasedTopicHastily(pbt *PullBasedTopic) []events.ReadOption {
	var output []events.ReadOption
	output = append(output, events.ReadLimit(uint(pbt.maxItems)))

	// output = append(output, events.ReadLimit(uint(pbt.maxItems)))

	if pbt.group != "" {
		output = append(output, natsjs.WithReadGroupName(pbt.group))
	}
	//we need to Ack manually
	output = append(output, natsjs.WithReadAutoAck(false))
	// natsjs.WithReadAutoAck(pbt.autoAck),
	return output
}

func (pbt *PullBasedTopic) read(ctx context.Context) ([]*events.Event, error) {
	pulledEvents, err := pbt.pullBasedConsumer.Read(
		pbt.topic, readOptionsFromPullBasedTopic(pbt)...,
	)
	if err != nil {
		return nil, err
	}
	return pulledEvents, nil
}

// TODO :: test process function.
func (pbt *PullBasedTopic) process(ctx context.Context, pulledEvents []*events.Event) error {
	//this is for delay configuration
	if pbt.readWithDelay {
		currentTime := time.Now()
		var readyToProcessEvents []*events.Event
		for _, event := range pulledEvents {
			if currentTime.After(event.Timestamp.Add(pbt.consumeDelay)) {
				readyToProcessEvents = append(readyToProcessEvents, event)
				if pbt.autoAck {
					// if the user sets auto ack we ack event immediately.
					//TODO :: handler error for auto ack .
					err := event.Ack()
					_ = err
				}
			} else {
				// we should nack the events that are not ready to be processed.
				err := event.Nack()
				_ = err
			}
		}
		if len(readyToProcessEvents) > 0 {

			err := pbt.handler(ctx, readyToProcessEvents)
			if err != nil {
				return err
			}
			if pbt.hastilyRecieveData && len(pulledEvents) == len(readyToProcessEvents) {
				// if we recieve maximum data we have to try it again hastily.
				pbt.hastilyRecieveDataNotif <- struct{}{}
			}
		}
		return nil
	}
	if err := pbt.handler(ctx, pulledEvents); err != nil {
		return err
	}
	if pbt.hastilyRecieveData && pbt.maxItems == len(pulledEvents) {
		pbt.hastilyRecieveDataNotif <- struct{}{}
	}
	return nil
}

func (pbt *PullBasedTopic) readHastily(ctx context.Context) ([]*events.Event, error) {
	readOpts := readOptionsFromPullBasedTopicHastily(pbt)

	pulledEvents, err := pbt.pullBasedConsumer.Read(
		pbt.topic, readOpts...,
	)
	if err != nil {
		return nil, err
	}
	return pulledEvents, err
}

func (pbt *PullBasedTopic) processHastily(ctx context.Context, pulledEvents []*events.Event) error {

	//this is for delay configuration
	if pbt.readWithDelay {
		currentTime := time.Now()
		var readyToProcessEvents []*events.Event
		for _, event := range pulledEvents {
			if currentTime.After(event.Timestamp.Add(pbt.consumeDelay)) {
				readyToProcessEvents = append(readyToProcessEvents, event)
				// if pbt.autoAck {
				// 	// if the user sets auto ack we ack event immediately.
				// 	//TODO :: handler error for auto ack .
				// 	event.Ack()
				// }
			} else {
				// we should nack the events that are not ready to be processed.
				err := event.Nack()
				_ = err
			}
		}
		if len(readyToProcessEvents) < pbt.minReqiredItmesRecieveHastily {
			for _, event := range readyToProcessEvents {
				err := event.Nack()
				_ = err
			}

			logger.Info("data pulled again hastily with consume delay but we need minReqiredItmesRecieveHastily to process events!")
			return nil
		} else if pbt.autoAck {
			for _, event := range readyToProcessEvents {
				err := event.Ack()
				_ = err
			}
		}
		err := pbt.handler(ctx, readyToProcessEvents)
		if err != nil {
			return err
		}
		if pbt.maxItems == len(readyToProcessEvents) {
			// in hastily check process if the size of readyToProcessEvents be equal to the maxItems, we repeate the process hastily again.
			pbt.hastilyRecieveDataNotif <- struct{}{}
		}
		return nil
	}
	if len(pulledEvents) < pbt.minReqiredItmesRecieveHastily {
		logger.Info("data pulled again hastily but we need minReqiredItmesRecieveHastily to process events!")

		for _, event := range pulledEvents {
			event.Nack()
		}

		return nil
	}

	if err := pbt.handler(ctx, pulledEvents); err != nil {
		return err
	} else if pbt.autoAck {
		for _, event := range pulledEvents {
			err = event.Ack()
			_ = err
		}
	}

	// if the buffer is full we have to hastily pull data from natsjs again.
	if pbt.maxItems == len(pulledEvents) {
		pbt.hastilyRecieveDataNotif <- struct{}{}
	}

	return nil
}
