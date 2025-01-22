package postgres

import "fmt"

type PostgresOptions struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	DBName     string `mapstructure:"dbName"`
	Password   string `mapstructure:"password"`
	Datasource string `mapstructure:"datasource"`
}

func (h *PostgresOptions) GetDatasource() string {
	datasource := fmt.Sprintf("postgres://%s/%s?sslmode=disable",
		h.Datasource,
		h.DBName,
	)

	return datasource
}

func (h *PostgresOptions) GetPostgresDatasource() string {
	datasource := fmt.Sprintf("postgres://%s/%s?sslmode=disable",
		h.Datasource,
		"postgres",
	)

	return datasource
}
