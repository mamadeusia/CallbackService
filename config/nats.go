package config

import (
	"strings"
	"time"
)

type Nats struct {
	URL      string
	Callback Callback
	Nkey     string
}

type Callback struct {
	Prefix   string
	Topics   string
	Final    string
	Duration string
}

type Max struct {
	Tx int
}

func NatsCallbackDurationString() string {
	switch cfg.Nats.Callback.Duration {
	case "Second":
		return "Second"
	case "Minute":
		return "Minute"
	case "Hour":
		return "Hour"
	default:
		return "Second"
	}
}

func NatsCallbackDuration() time.Duration {
	switch cfg.Nats.Callback.Duration {
	case "Second":
		return time.Second
	case "Minute":
		return time.Minute
	case "Hour":
		return time.Hour
	default:
		return time.Second
	}
}

func NatsCallbackTopics() []string {
	if cfg.Nats.Callback.Topics == "" {
		return []string{"10", "15", "30", "60", "120"}
	}
	return strings.Split(cfg.Nats.Callback.Topics, ",")
}

func NatsCallbackPrefix() string {
	if cfg.Nats.Callback.Prefix == "" {
		return "CALLBACKS"
	}
	return cfg.Nats.Callback.Prefix
}

func NatsCallbackFinalTopic() string {
	if cfg.Nats.Callback.Final == "" {
		return "final"
	}
	return cfg.Nats.Callback.Final
}

func NatsNkey() string {
	return cfg.Nats.Nkey
}
func NatsURL() string {
	if cfg.Nats.URL == "" {
		return "localhost:4222"
	}
	return cfg.Nats.URL
}
