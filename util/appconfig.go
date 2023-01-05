package util

import "github.com/spf13/viper"
import log "github.com/sirupsen/logrus"

type AppSettings struct {
	Port     string   `mapstructure:"port"`
	LogLevel string   `mapstructure:"logLevel"`
	Regions  []string `mapstructure:"regions"`
}

func LoadAppConfig() AppSettings {
	var settings AppSettings
	log.Println("Loading Server Configurations...")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&settings)
	if err != nil {
		log.Fatal(err)
	}
	return settings

}
