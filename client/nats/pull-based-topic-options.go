package nats

import "time"

func WithPullBasedGroup(group string) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.group = group
		return nil
	}
}

func WithPullBasedTopicString(topic string) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.topic = topic
		return nil
	}
}

// duration specifies the duration time between each pulls from the topic.
func WithPullBasedDuration(duration time.Duration) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.duration = duration
		return nil
	}
}

// maxItems is the maximum number of items that we can pull every time we pull from the	topic.
func WithPullBasedMaxItems(maxItems int) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.maxItems = maxItems
		return nil
	}
}

// eventHandler is the function that handles events, you can create your own event handler.
func WithPullBasedEventHandler(eventHandler EventHandler) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.handler = eventHandler
		return nil
	}
}

// if it's false then you have to ack the event in the event handler manually.
func WithPullBasedAutoAck(autoAck bool) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.autoAck = autoAck
		return nil
	}
}

// if we need to create delayed and topics WithPullBased delayed we need to set this
func WithPullBasedConsumeDelayTime(delayTime time.Duration) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.readWithDelay = true
		p.consumeDelay = delayTime
		return nil
	}
}

// if you set this option, we try to pull data hastily if the maxItems is reached,
func WithPullBasedHastilyRecieveData(minItmesRecievedHastily int) PullBasedTopicConfiguration {
	return func(p *PullBasedTopic) error {
		p.hastilyRecieveData = true
		p.minReqiredItmesRecieveHastily = minItmesRecievedHastily
		return nil
	}
}
