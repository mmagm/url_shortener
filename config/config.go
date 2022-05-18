package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	ListenPort string `json:"listenPort"`
}

type DatabaseConfiguration struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	SSLMode  string `json:"sslmode"`
}

type Configuration struct {
	Server      ServerConfiguration   `json:"server"`
	Database    DatabaseConfiguration `json:"database"`
	Environment string                `json:"environment"`
	BaseUrl     string                `json:"baseUrl"`
	ApiToken    string                `json:"apiToken"`
}

func Load() (c Configuration, err error) {
	viper.SetConfigName("config") // name of config file (without extension)

	if os.Getenv("CONFIG_PATH") != "" {
		viper.AddConfigPath(os.Getenv("CONFIG_PATH"))
	}
	viper.AddConfigPath(".")                             // optionally look for config in the working directory
	viper.AddConfigPath("./test_configuration")          // optionally look for config in the working directory

	err = viper.ReadInConfig()
	if err != nil {
		return c, fmt.Errorf("Fatal error config file: %s \n", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return c, fmt.Errorf("Fatal error config file: %s \n", err)
	}
	return
}
