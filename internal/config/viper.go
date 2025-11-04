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

	if err := godotenv.Load(".env"); err != nil {
		log.Warn(".env file not found or unreadable")
	} else {
		log.Info(".env loaded via godotenv")
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v
}
