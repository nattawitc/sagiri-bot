package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("sagiri")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetConfigType("yaml")
	viper.SetConfigName("main")
	viper.AddConfigPath("conf/")
	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("config file not found")
		default:
			panic(err)
		}
	}
}
