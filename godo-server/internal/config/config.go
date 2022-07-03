package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseHost 	 string `mapstructure:"DB_HOST"`
	DatabasePort 	 string `mapstructure:"DB_PORT"`
	DatabaseName 	 string `mapstructure:"DB_NAME"`
	DatabaseUsername string `mapstructure:"DB_UNAME"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	ApiPort		     string `mapstructure:"API_PORT"`
}

func LoadDevConfig() (conf Config) {
	return makeConfig("dev")
}

func LoadConfig() (conf Config)  {
	return makeConfig("config")
}

func makeConfig(configName string) (conf Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		log.Fatal("There has been an issue reading the configuration file.")
		return
	}
	
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Println(err)
		log.Fatal("There has been an issue un-marshaling the configuration file.")
		return
	}

	return
}