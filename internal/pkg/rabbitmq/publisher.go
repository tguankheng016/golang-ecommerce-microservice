package rabbitmq

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/gofrs/uuid"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
)

type IPublisher interface {
	PublishMessage(msg interface{}) error
	IsPublished(msg interface{}) bool
}

var publishedMessages []string

// TODO: Add Jaeger Tracer OpenTelemetry
type Publisher struct {
	cfg  *RabbitMQOptions
	conn *amqp.Connection
	//jaegerTracer trace.Tracer
	ctx context.Context
}

func NewPublisher(ctx context.Context, cfg *RabbitMQOptions, conn *amqp.Connection) IPublisher {
	return &Publisher{ctx: ctx, cfg: cfg, conn: conn}
}

func (p Publisher) PublishMessage(msg interface{}) error {

	data, err := jsoniter.Marshal(msg)

	if err != nil {
		logger.Logger.Error("Error in marshalling message to publish message")
		return err
	}

	typeName := reflect.TypeOf(msg).Elem().Name()
	snakeTypeName := strcase.ToSnake(typeName)

	// Temp
	ctx := p.ctx
	// ctx, span := p.jaegerTracer.Start(p.ctx, typeName)
	// defer span.End()

	// Inject the context in the headers
	//headers := otel.InjectAMQPHeaders(ctx)

	channel, err := p.conn.Channel()
	if err != nil {
		logger.Logger.Error("Error in opening channel to consume message")
		return err
	}

	defer channel.Close()

	err = channel.ExchangeDeclare(
		snakeTypeName, // name
		p.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		logger.Logger.Error("Error in declaring exchange to publish message")
		return err
	}

	correlationId := ""

	if ctx.Value(echo.HeaderXCorrelationID) != nil {
		correlationId = ctx.Value(echo.HeaderXCorrelationID).(string)
	}

	messageId, err := uuid.NewV4()
	if err != nil {
		logger.Logger.Error("Error in generating message uuid")
		return err
	}
	publishingMsg := amqp.Publishing{
		Body:          data,
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		MessageId:     messageId.String(),
		Timestamp:     time.Now(),
		CorrelationId: correlationId,
		//Headers:       headers,
	}

	err = channel.Publish(snakeTypeName, snakeTypeName, false, false, publishingMsg)

	if err != nil {
		logger.Logger.Error("Error in publishing message")
		return err
	}

	publishedMessages = append(publishedMessages, snakeTypeName)

	// h, err := jsoniter.Marshal(headers)
	// if err != nil {
	// 	logger.Logger.Error("Error in marshalling headers to publish message")
	// 	return err
	// }

	logger.Logger.Info(fmt.Sprintf("Published message: %s", publishingMsg.Body))
	// span.SetAttributes(attribute.Key("message-id").String(publishingMsg.MessageId))
	// span.SetAttributes(attribute.Key("correlation-id").String(publishingMsg.CorrelationId))
	// span.SetAttributes(attribute.Key("exchange").String(snakeTypeName))
	// span.SetAttributes(attribute.Key("kind").String(p.cfg.Kind))
	// span.SetAttributes(attribute.Key("content-type").String("application/json"))
	// span.SetAttributes(attribute.Key("timestamp").String(publishingMsg.Timestamp.String()))
	// span.SetAttributes(attribute.Key("body").String(string(publishingMsg.Body)))
	// span.SetAttributes(attribute.Key("headers").String(string(h)))

	return nil
}

func (p *Publisher) IsPublished(msg interface{}) bool {

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)
	isPublished := slices.Contains(publishedMessages, snakeTypeName)

	return isPublished
}
