package otel

type JaegerOptions struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	ServiceName string `mapstructure:"serviceName"`
	TracerName  string `mapstructure:"tracerName"`
}
