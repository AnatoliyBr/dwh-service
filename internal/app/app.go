package app

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnatoliyBr/dwh-service/internal/controller/apiserver"
	"github.com/AnatoliyBr/dwh-service/internal/repository"
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func Run() {

	// PostgreSQL
	if err := godotenv.Load(); err != nil {
		logrus.Fatal(fmt.Errorf("app - Run - godotenv.Load: %w", err))
	}
	configDB := repository.NewConfig()
	db, err := repository.NewDB(configDB)
	if err != nil {
		logrus.Fatal(fmt.Errorf("app - Run - repository.NewDB: %w", err))
	}
	defer db.Close()

	// Controller
	flag.Parse()
	configAPIServer := apiserver.NewConfig()
	_, err = toml.DecodeFile(configPath, configAPIServer)
	if err != nil {
		logrus.Fatal(fmt.Errorf("app - Run - toml.DecodeFile: %w", err))
	}

	s, err := apiserver.NewAPIServer(configAPIServer)
	if err != nil {
		logrus.Fatal(fmt.Errorf("app - Run - apiServer.NewAPIServer: %w", err))
	}

	s.StartAPIServer()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case signal := <-interrupt:
		logrus.Info("app - Run - signal: " + signal.String())
	case err = <-s.Notify():
		logrus.Error(fmt.Errorf("app - Run - apiServer.Notify: %w", err))
	}

	// Shutdown
	err = s.Shutdown()
	if err != nil {
		logrus.Error(fmt.Errorf("app - Run - apiServer.Shutdown: %w", err))
	}
}
