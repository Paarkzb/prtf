package main

import (
	"os"
	"os/signal"
	"prtf"
	"prtf/pkg/handler"
	"prtf/pkg/repository"
	"prtf/pkg/service"
	"syscall"

	"golang.org/x/net/context"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}

	dbconfig := repository.DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
	db, err := repository.NewPostgresDB(dbconfig)
	if err != nil {
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(db)
	serivces := service.NewService(repos)
	handlers := handler.NewHandler(serivces)

	srv := new(prtf.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Printf("server listening on port %s", viper.GetString("port"))

	logrus.Print("prtf-server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("prtf-server shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	db.Close()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	return viper.ReadInConfig()
}
