package config

import (
	"os"
	"pura-agung-kertajaya-backend/internal/util"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(viper.GetInt32("log.level")))
	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetOutput(os.Stdout)

	token := viper.GetString("BETTERSTACK_TOKEN")
	if token != "" {
		log.AddHook(&util.BetterStackHook{
			Token: token,
			URL:   "https://in.logs.betterstack.com",
		})
	}

	return log
}
