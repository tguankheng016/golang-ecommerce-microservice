package messaging

type WatermillNatsOptions struct {
	Enabled     bool   `mapstructure:"enabled"`
	Url         string `mapstructure:"url"`
	DurableName string `mapstructure:"durableName"`
}

type WatermillOptions struct {
	Nats *WatermillNatsOptions `mapstructure:"nats"`
}
