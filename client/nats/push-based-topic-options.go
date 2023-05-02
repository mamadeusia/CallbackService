package nats

import "time"

func WithPushBasedTopicString(topic string) PushBasedTopicConfiguration {
	return func(p *PushBasedTopic) error {
		p.topic = topic
		return nil
	}
}

// // maxItems is the maximum number of items that we can pull every time we pull from the	topic.
// func WithMaxItems(maxItems int) PullBasedTopicConfiguration {
// 	return func(p *PullBasedTopic) error {
// 		p.maxItems = maxItems
// 		return nil
// 	}
// }

// eventHandler is the function that handles events, you can create your own event handler.
func WithPushBasedEventHandler(eventHandler EventHandler) PushBasedTopicConfiguration {
	return func(p *PushBasedTopic) error {
		p.handler = eventHandler
		return nil
	}
}

// if it's false then you have to ack the event in the event handler manually.
func WithPushBasedAutoAck(autoAck bool) PushBasedTopicConfiguration {
	return func(p *PushBasedTopic) error {
		p.autoAck = autoAck
		return nil
	}
}

// if we need to create delayed and topics with delayed we need to set this
func WithPushBasedConsumeDelayTime(delayTime time.Duration) PushBasedTopicConfiguration {
	return func(p *PushBasedTopic) error {
		p.readWithDelay = true
		p.consumeDelay = delayTime
		return nil
	}
}

func WithPushBasedBatchSize(batchSize int) PushBasedTopicConfiguration {
	return func(p *PushBasedTopic) error {
		p.batchEnabled = true
		p.batchSize = batchSize
		return nil
	}
}

// // if you set this option, we try to pull data hastily if the maxItems is reached,
// func WithHastilyRecieveData(minItmesRecievedHastily int) PullBasedTopicConfiguration {
// 	return func(p *PullBasedTopic) error {
// 		p.hastilyRecieveData = true
// 		p.minItmesRecievedHastily = minItmesRecievedHastily
// 		return nil
// 	}
// }
