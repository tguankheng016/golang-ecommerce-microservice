package http

import "fmt"

type ServerOptions struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Name        string `mapstructure:"name"`
	CorsOrigins string `mapstructure:"corsOrigins"`
}

func (h *ServerOptions) GetBasePath() string {
	basePath := fmt.Sprintf("http://%s:%s", h.Host, h.Port)

	return basePath
}
