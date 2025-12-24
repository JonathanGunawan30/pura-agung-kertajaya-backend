package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	log := logrus.New()
	v := viper.New()

	err := godotenv.Load(".env")

	if err != nil {
		err = godotenv.Load("../.env")
	}

	if err != nil {
		err = godotenv.Load("../../.env")
	}

	if err != nil {
		log.Warn(".env file not found in current or parent directory, utilizing system environment variables")
	} else {
		log.Info(".env loaded successfully via godotenv")
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v
}
