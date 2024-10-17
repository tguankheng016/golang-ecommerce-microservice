package rabbitmq

type RabbitMQOptions struct {
	Host         string
	Port         int
	User         string
	Password     string
	ExchangeName string
	Kind         string
}
