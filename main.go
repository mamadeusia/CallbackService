package main

import (
	"context"

	"github.com/go-micro/plugins/v4/events/natsjs"

	"github.com/mamadeusia/CallbackService/config"
	"github.com/mamadeusia/CallbackService/entity"
	"github.com/mamadeusia/CallbackService/handler/callbackevent"
	"github.com/mamadeusia/CallbackService/service/callback"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	grpcs "github.com/go-micro/plugins/v4/server/grpc"
)

var (
	service = "connector-callback-service"
	version = "latest"
)

func main() {
	if err := config.Load(); err != nil {
		logger.Fatal("could not load config:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	srv := micro.NewService(
		micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
		micro.Context(ctx),

		micro.BeforeStop(func() error {
			cancel()
			return nil
		}),
		// micro.AfterStop()
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
	)

	var natsjsConfigOptions []natsjs.Option
	natsjsConfigOptions = append(natsjsConfigOptions, natsjs.Address(config.NatsURL()))

	if config.NatsNkey() != "" {
		natsjsConfigOptions = append(natsjsConfigOptions, natsjs.NkeyConfig(config.NatsNkey()))
	}

	stream, err := natsjs.NewStream(natsjsConfigOptions...)
	if err != nil {
		logger.Fatal(err)
	}

	store, err := natsjs.NewStore(natsjsConfigOptions...)

	if err != nil {
		logger.Error(err)
	}

	callbackSrv, err := callback.NewService(
		callback.WithFailureStream(stream),
	)
	if err != nil {
		logger.Fatal(err)
	}

	handler, err := callbackevent.NewHandler(ctx, config.NatsCallbackPrefix(), config.NatsCallbackTopics(), callbackSrv, store, stream)
	if err != nil {
		logger.Fatal(err)
	}

	if err := handler.StartConsumer(ctx); err != nil {
		logger.Fatal(err)
	}

	stream.Publish("CALLBACKS.10Second", entity.CallBackData{
		Url:  "https://webhook.site/77c3a62c-d416-425a-8817-8d46194f8a44",
		Data: []byte("First Message"),
	})

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
