package http

import "fmt"

type ServerOptions struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Name        string `mapstructure:"name"`
	AppUrl      string `mapstructure:"appUrl"`
	CorsOrigins string `mapstructure:"corsOrigins"`
}

func (h *ServerOptions) GetBasePath() string {
	basePath := h.AppUrl

	if basePath == "" {
		basePath = fmt.Sprintf("http://%s:%s", h.Host, h.Port)
	}

	return basePath
}
