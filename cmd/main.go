package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testTask"
	"testTask/handler"
	"testTask/repository"
	"testTask/service"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("config initialization error: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(testTask.Server)
	go startServer(server, handlers)
	waitForShutdown(server, db)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return viper.ReadInConfig()
}

func startServer(server *testTask.Server, handlers *handler.Handler) {
	if err := server.Run(viper.GetString("server.port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}
	logrus.Printf("Server started on port %s", viper.GetString("port"))
}

func waitForShutdown(server *testTask.Server, db *sqlx.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}
}
