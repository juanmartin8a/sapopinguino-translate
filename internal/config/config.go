package config

// READS THE .yml FILES FROM ./config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	c config
}

type config struct {
	Websocket struct {
		Endpoint string `mapstructure:"endpoint"`
	} `mapstructure:"websocket"`
	OpenAI struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"openai"`
}

func LoadConfig() (*Config, error) {

	setEnvConfig()

	viper.SetConfigType("yml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error while running viper.ReadInConfig(): %v", err)
	}

	var tempConfig config

	if err := viper.Unmarshal(&tempConfig); err != nil {
		return nil, fmt.Errorf("Error while running viper.Unmarshal(): %v", err)
	}

	c := &Config{
		c: tempConfig,
	}

	return c, nil
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	dir := path.Dir(b)
	return filepath.Join(dir, "../..") // sapopinguino-translate/
}

func setEnvConfig() {
	viper.AddConfigPath(
		filepath.Join(RootDir(), "config"),
	)
	viper.SetConfigName(config_filename)
}

func (c *Config) OpenAIKey() string {
	return c.c.OpenAI.Key
}

func (c *Config) WebsocketEndpoint() *string { // pointer to string because that's what the SDK asks for
	return &c.c.Websocket.Endpoint
}
