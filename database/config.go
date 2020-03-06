package database

type Config struct {
	Driver         string `toml:"driver"`
	DataSourceName string `toml:"data-source-name"`
	CreateScript   string `toml:"create-script"`
	CreateQuery string `toml:"create-query"`
}

