package callbackevent

import (
	"context"
	"strconv"
	"time"

	"github.com/mamadeusia/CallbackService/client/nats"
	"github.com/mamadeusia/CallbackService/config"
	"github.com/mamadeusia/CallbackService/service/callback"
	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

const (
	groupName = "callback"
)

type Handler struct {
	callbackSrv callback.Callback
	topics      []string
	prefixTopic string

	store  events.Store
	stream events.Stream
}

func NewHandler(ctx context.Context,
	prefixTopic string,
	topics []string,
	callbackService callback.Callback,
	store events.Store,
	stream events.Stream) (*Handler, error) {
	return &Handler{
		callbackSrv: callbackService,
		topics:      topics,
		prefixTopic: prefixTopic,
		store:       store,
		stream:      stream,
	}, nil

}

func (h *Handler) StartConsumer(ctx context.Context) error {
	//shared configuration between pull based topics.
	sharedPullBasedConfigurations := []nats.PullBasedTopicConfiguration{
		nats.WithPullBasedGroup(groupName),
		nats.WithPullBasedAutoAck(true),
	}
	// configure pull based configurations.
	var pullBasedTopics []*nats.PullBasedTopic

	for i, topic := range h.topics {

		topicConfiguration := append(sharedPullBasedConfigurations, nats.WithPullBasedTopicString(config.NatsCallbackPrefix()+"."+topic+config.NatsCallbackDurationString()))
		topicConsumeDelay, err := strconv.Atoi(topic)
		if err != nil {
			logger.Fatal(err)
		}
		topicConfiguration = append(topicConfiguration, nats.WithPullBasedMaxItems(10))

		topicConfiguration = append(topicConfiguration, nats.WithPullBasedDuration(3*time.Second))

		topicConfiguration = append(topicConfiguration, nats.WithPullBasedConsumeDelayTime(time.Duration(topicConsumeDelay*int(config.NatsCallbackDuration()))))

		if i == len(h.topics)-1 {
			topicConfiguration = append(topicConfiguration, nats.WithPullBasedEventHandler(
				nats.EventHandler(h.CallWithRepublish_FailureScenario(config.NatsCallbackPrefix()+"."+config.NatsCallbackFinalTopic()))),
			)
		} else {
			nextTopic := h.topics[i+1]
			topicConfiguration = append(topicConfiguration, nats.WithPullBasedEventHandler(
				nats.EventHandler(h.CallWithRepublish_FailureScenario(config.NatsCallbackPrefix()+"."+nextTopic+config.NatsCallbackDurationString()))),
			)
		}

		newTopic, err := nats.NewPullBasedTopic(topicConfiguration...)
		if err != nil {
			logger.Fatal(err)
		}
		pullBasedTopics = append(pullBasedTopics, newTopic)
	}

	var natConfigurations []nats.NatsConfiguration

	natConfigurations = append(natConfigurations, nats.WithPullBasedStore(h.store))
	natConfigurations = append(natConfigurations, nats.WithPushBasedStream(h.stream))

	for _, pullBasedTopic := range pullBasedTopics {
		natConfigurations = append(natConfigurations, nats.WithPullBasedTopic(pullBasedTopic))
	}

	natsClient, err := nats.New(natConfigurations...)
	if err != nil {
		logger.Fatal(err)
	}

	return natsClient.Start(ctx)

}
