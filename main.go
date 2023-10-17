package main

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/jcpribeiro/TransactionApp/config"
	"github.com/jcpribeiro/TransactionApp/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadConfig(path string) error {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&config.GlobalConfig)
}

func main() {
	if os.Getenv("APP_ENV") == "prod"{
		err := loadConfig("./config_prod.json")
		if err != nil {
			logrus.Fatal("cannot load config: ", err)
		}
	}else {
		err := loadConfig("./config.json")
		if err != nil {
			logrus.Fatal("cannot load config: ", err)
		}
	}

	server := server.NewServer()
	go server.Start()
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	logrus.Info("Stopped with signal: ", <-c)
}
