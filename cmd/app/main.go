package main

import (
	"fmt"
	"log"

	"github.com/KatenkaKet/wallet"
	"github.com/KatenkaKet/wallet/pkg/handler"
	"github.com/KatenkaKet/wallet/pkg/repository"
	"github.com/KatenkaKet/wallet/pkg/service"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

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

	fmt.Println(viper.GetString("PORT"))

	srv := new(wallet.Server)
	if err := srv.Run(viper.GetString("PORT"), hdl.InitRoutes()); err != nil {
		log.Fatal("error occured while running http server: ", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	return viper.ReadInConfig()
}
