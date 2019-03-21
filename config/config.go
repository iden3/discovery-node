package config

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type Config struct {
	Datadir string
	Ports   struct {
		WebSockets int
		HTTPRPC    int
		Bzz        int
		P2P        int
	}
	Pss struct {
		Kind     string
		Key      string
		Topic    string
		LogLevel string
	}
}

var C Config

func MustRead(c *cli.Context) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if c.GlobalString("config") != "" {
		viper.SetConfigFile(c.GlobalString("config"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	return nil
}
