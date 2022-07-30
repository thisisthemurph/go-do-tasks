package configuration

import (
	"godo/internal/helper/ilog"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseHost     string `mapstructure:"DB_HOST"`
	DatabasePort     string `mapstructure:"DB_PORT"`
	DatabaseName     string `mapstructure:"DB_NAME"`
	DatabaseUsername string `mapstructure:"DB_UNAME"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	ApiPort          string `mapstructure:"API_PORT"`
	JWTKey           string `mapstructure:"JWT_KEY"`
}

func LoadDevConfig(logger ilog.StdLogger) (conf Config) {
	return makeConfig("dev", logger)
}

func LoadConfig(logger ilog.StdLogger) (conf Config) {
	return makeConfig("config", logger)
}

func makeConfig(configName string, log ilog.StdLogger) (conf Config) {
	viper.SetConfigName(configName)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorln(err)
		log.Fatal("There has been an issue reading the configuration file.")
		return
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Errorln(err)
		log.Fatal("There has been an issue un-marshaling the configuration file.")
		return
	}

	return
}
