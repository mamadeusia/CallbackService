package callbackevent

import (
	"context"
	"encoding/json"

	"github.com/mamadeusia/CallbackService/client/nats"
	"github.com/mamadeusia/CallbackService/entity"
	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

func (h *Handler) CallWithRepublish_FailureScenario(nextTopic string) nats.EventHandler {
	return func(ctx context.Context, e []*events.Event) error {
		for _, event := range e {

			callbackData := entity.CallBackData{}
			err := json.Unmarshal(event.Payload, &callbackData)
			if err != nil {
				logger.Error("call with republish failed unmarshal ,", err)
				continue
			}

			err = h.callbackSrv.SendDataCallBack(ctx, callbackData, nextTopic)
			if err != nil {
				logger.Error("call with republish failed to send data to callback service, error: %v", err)
			}

			// if err := event.Ack(); err != nil {
			// 	logger.Error("call with republish fail to ack event, error: %v", err)
			// }

		}
		return nil
	}
}
