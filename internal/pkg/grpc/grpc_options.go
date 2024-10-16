package grpc

type GrpcOptions struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Development bool   `mapstructure:"development"`
	Enabled     bool   `mapstructure:"enabled"`
}
