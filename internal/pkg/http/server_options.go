package http

type ServerOptions struct {
	Port        string `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Name        string `mapstructure:"name"`
	AppUrl      string `mapstructure:"appUrl"`
	CorsOrigins string `mapstructure:"corsOrigins"`
}

func (h *ServerOptions) GetBasePath() string {
	basePath := h.AppUrl

	return basePath
}
