package messaging

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	nc "github.com/nats-io/nats.go"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
)

func NewWatermillPublisher(env environment.Environment, logger watermill.LoggerAdapter, config *WatermillOptions) (message.Publisher, error) {
	if env.IsTest() || !config.Nats.Enabled {
		// Use In Memory Publisher
		pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)
		return pubSub, nil
	}

	natsURL, marshaler, options, jsConfig := GetNatsConfig(config)

	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         natsURL,
			NatsOptions: options,
			Marshaler:   marshaler,
			JetStream:   jsConfig,
		},
		logger,
	)

	return publisher, err
}

func GetNatsConfig(config *WatermillOptions) (string, *nats.GobMarshaler, []nc.Option, nats.JetStreamConfig) {
	natsURL := config.Nats.Url
	marshaler := &nats.GobMarshaler{}

	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
	}

	subscribeOptions := []nc.SubOpt{
		nc.DeliverAll(),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}

	return natsURL, marshaler, options, jsConfig
}
