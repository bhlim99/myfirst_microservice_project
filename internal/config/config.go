package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	DBConfig struct {
		Host         string `mapstructure:"PG_HOST"`
		Port         int    `mapstructure:"PG_PORT_NUMBER"`
		Username     string `mapstructure:"PG_USERNAME"`
		Password     string `mapstructure:"PG_PASSWORD"`
		DatabaseName string `mapstructure:"PG_NAME"`
		SslMode      string `mapstructure:"PG_SSL_MODE"`
		Source       string
		SourceUrl    string
		// DefaultMaxOpenConnections the default value for max open connections in the PostgreSQL connection pool
		MaxOpenConnections     int    `mapstructure:"PG_MAX_OPEN_CONNECTIONS"`
		MaxIdleConnections     int    `mapstructure:"PG_MAX_IDLE_CONNECTIONS"`
		MigrationFileUrl       string `mapstructure:"PG_MIGRATION_FILE_URL"`
		ConnectionsMaxLifetime time.Duration
	}
	app struct {
		AppName    string `mapstructure:"APP_NAME"`
		AppEnv     string `mapstructure:"APP_ENV"`
		AppVersion string `mapstructure:"APP_VERSION"`
		DB         *DBConfig
		// Http       *HttpConfig
		// Grpc       *GrpcConfig
	}

	Config struct {
		App *app
	}
)

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.App); err != nil {
		return
	}
	if err = viper.Unmarshal(&config.App.DB); err != nil {
		return
	}
	config.App.DB.SourceUrl = fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.App.DB.Username,
		config.App.DB.Password,
		config.App.DB.Host,
		config.App.DB.Port,
		config.App.DB.DatabaseName,
		config.App.DB.SslMode,
	)

	config.App.DB.Source = fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		config.App.DB.Host,
		config.App.DB.Port,
		config.App.DB.Username,
		config.App.DB.DatabaseName,
		config.App.DB.Password,
		config.App.DB.SslMode,
	)
	return
}
