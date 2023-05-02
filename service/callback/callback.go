package callback

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/mamadeusia/CallbackService/entity"
	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

type ServiceConfiguration func(s *Service) error

type EventHandler func(ctx context.Context, e []*events.Event) error

type Callback interface {
	SendDataCallBack(ctx context.Context, data entity.CallBackData, nextTopic string) error
}

type Service struct {
	Stream events.Stream
}

func NewService(cfgs ...ServiceConfiguration) (*Service, error) {
	service := &Service{}
	for _, cfg := range cfgs {
		if err := cfg(service); err != nil {
			return nil, err
		}
	}
	return service, nil
}

func (s *Service) SendDataCallBack(ctx context.Context, data entity.CallBackData, nextTopic string) error {

	logger.Info("Call Back data : ", nextTopic)
	bodyReader := bytes.NewReader(data.Data)
	rsp, err := http.Post(data.Url, "multipart/form-data", bodyReader)
	if err != nil || rsp.StatusCode > 300 || rsp.StatusCode < 200 {
		err = errors.New("can't post data")

		if err := s.Stream.Publish(nextTopic, data); err != nil {
			logger.Error(err)
			return err
		}
		return err

	}

	return nil

}
