package messaging

type WatermillNatsOptions struct {
	Enabled bool   `mapstructure:"enabled"`
	Url     string `mapstructure:"url"`
}

type WatermillOptions struct {
	Nats *WatermillNatsOptions `mapstructure:"nats"`
}
