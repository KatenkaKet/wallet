package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KatenkaKet/wallet"
	"github.com/KatenkaKet/wallet/pkg/handler"
	"github.com/KatenkaKet/wallet/pkg/repository"
	"github.com/KatenkaKet/wallet/pkg/service"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// Запускать из корня проекта!

func main() {
	if err := initConfig(); err != nil {
		log.Fatal("error initializing config: ", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Username: viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSLMODE"),
	})

	if err != nil {
		log.Fatal("error initializing database: ", err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	hdl := handler.NewHandler(service)

	//fmt.Println(viper.GetString("PORT"))

	srv := new(wallet.Server)

	go func() {
		if err := srv.Run(viper.GetString("PORT"), hdl.InitRoutes()); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal("error occurred while running http server: ", err.Error())
			}
		}
	}()

	log.Println("Listening on " + viper.GetString("PORT"))

	quet := make(chan os.Signal, 1)
	signal.Notify(quet, syscall.SIGINT, syscall.SIGTERM)
	<-quet

	log.Println("Shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("error occured while shutting down http server: ", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatal("error occured while closing database: ", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	return viper.ReadInConfig()
}
