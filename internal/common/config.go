package common

import (
	"flag"
	"strings"

	"github.com/dmitrii-a/image-previewer/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// CacheConfig cache config.
type CacheConfig struct {
	MaxSize int64 `mapstructure:"MAX_SIZE"`
}

// ServerConfig server config.
type ServerConfig struct {
	Host            string `mapstructure:"HOST"`
	Port            int    `mapstructure:"PORT"`
	Debug           bool   `mapstructure:"DEBUG"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	ShutdownTimeout int    `mapstructure:"SHUTDOWN_TIMEOUT_SECOND"`
	ReadTimeout     int    `mapstructure:"READ_TIMEOUT_SECOND"`
	WriteTimeout    int    `mapstructure:"WRITE_TIMEOUT_SECOND"`
}

// AppConfig app config.
type AppConfig struct {
	Server ServerConfig `mapstructure:"APP"`
	Cache  CacheConfig  `mapstructure:"CACHE"`
}

// Config project config.
var Config AppConfig

func setDefaults() {
	viper.SetDefault("APP.HOST", "127.0.0.1")
	viper.SetDefault("APP.PORT", 8080)
	viper.SetDefault("APP.DEBUG", true)
	viper.SetDefault("APP.LOG_LEVEL", "info")
	viper.SetDefault("APP.SHUTDOWN_TIMEOUT_SECOND", 30)
	viper.SetDefault("APP.READ_TIMEOUT_SECOND", 10)
	viper.SetDefault("APP.WRITE_TIMEOUT_SECOND", 10)

	viper.SetDefault("CACHE.MAX_SIZE", 1024*1024*100)
}

func init() {
	var err error

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	log := logger.InitLogger(zerolog.GlobalLevel())

	setDefaults()

	if IsErr(err) {
		log.Fatal().Msgf("Unable to decode into struct, %v", err)
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&Config)
	if IsErr(err) {
		log.Fatal().Msgf("Unable to decode into struct, %v", err)
	}
}

// SetConfigFileSettings applying settings from a configuration file.
func (config *AppConfig) SetConfigFileSettings(path string) {
	flag.Parse()

	log := logger.InitLogger(zerolog.GlobalLevel())

	if path != "" {
		viper.SetConfigFile(path)

		if err := viper.ReadInConfig(); IsErr(err) {
			log.Fatal().Msgf("Error reading config file, %s", err)
		}

		err := viper.Unmarshal(&Config)
		if IsErr(err) {
			log.Fatal().Msgf("Unable to decode into struct, %v", err)
		}
	}
	// Override config with environment variables if they exist.
	viper.AutomaticEnv()

	err := setLogger(Config.Server.LogLevel)
	if IsErr(err) {
		log.Fatal().Msgf("Error parsing loglevel, %s", err)
	}
}
