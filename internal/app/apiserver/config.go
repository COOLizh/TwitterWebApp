package apiserver

// Config ...
type Config struct {
	BindAddr           string `toml:"bind_addr"`
	LogLevel           string `toml:"log_level"`
	DbConnectionString string `toml:"db_connection_string"`
	DbName             string `toml:"db_name"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr:           ":8080",
		LogLevel:           "debug",
		DbConnectionString: "mongodb://localhost:27017",
		DbName:             "twitter_web_app",
	}
}
