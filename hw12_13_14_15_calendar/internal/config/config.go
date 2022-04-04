package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	Dir   string `mapstructure:"dir"`
}

type StorageConf struct {
	Type     string `mapstructure:"type"`
	Login    string `mapstructure:"login"`
	Password string `mapstructure:"password"`
}

type ServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func NewConfig(configFile string) Config {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	c := &Config{
		Logger: LoggerConf{
			Level: "DEBUG",
			Dir:   "logs",
		},
		Storage: StorageConf{
			Type:     "MEMO",
			Login:    "",
			Password: "",
		},
		Server: ServerConf{
			Host: "127.0.0.1",
			Port: "8888",
		},
	}
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("couldn't read config: %s", err)
	}
	return *c
}
