package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	log := logrus.New()
	config := viper.New()

	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Warnf("Failed to load .env file from current or parent directory: %v", err)
		}
	}

	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	config.AddConfigPath("./config")
	config.AddConfigPath("/")

	err = config.ReadInConfig()

	if err != nil {
		log.Fatalf("Fatal error config file: %v \n", err)
	}

	log.Infof("Configuration loaded successfully")
	return config
}
