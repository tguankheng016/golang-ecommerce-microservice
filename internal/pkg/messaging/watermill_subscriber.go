package messaging

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
)

func NewWatermillSubscriber(env environment.Environment, logger watermill.LoggerAdapter, config *WatermillOptions) (message.Subscriber, error) {
	if env.IsTest() || !config.Nats.Enabled {
		// Use In Memory Subscriber
		pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)
		return pubSub, nil
	}

	natsURL, marshaler, options, jsConfig := GetNatsConfig(config)

	subscriber, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:            natsURL,
			CloseTimeout:   30 * time.Second,
			AckWaitTimeout: 30 * time.Second,
			NatsOptions:    options,
			Unmarshaler:    marshaler,
			JetStream:      jsConfig,
		},
		logger,
	)

	return subscriber, err
}
