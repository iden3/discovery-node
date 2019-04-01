package config

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Config holds the config file data
type Config struct {
	Datadir  string
	DbPath   string
	KeyStore struct {
		Path     string
		Password string
	}
	Ports struct {
		API        int
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
	DiscoverFreshTimeout int64
	Mode                 string
}

// C contains the Config data
var C Config

// MustRead reads the config file and puts the data into the C variable
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
