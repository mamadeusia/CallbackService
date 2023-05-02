package callback

import (
	"go-micro.dev/v4/events"
)

func WithFailureStream(stream events.Stream) ServiceConfiguration {
	return func(s *Service) error {
		s.Stream = stream
		return nil
	}
}
