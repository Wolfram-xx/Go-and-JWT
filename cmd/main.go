package main

import (
	"REST_JWT"
	handler2 "REST_JWT/package/handler"
	"REST_JWT/package/repository"
	"REST_JWT/package/service"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if errs := initConfig(); errs != nil {
		log.Fatalf("ERROR of configurations: %s", errs.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Username: viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASS"),
		Database: viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("ERROR of connection to DB: %s", err.Error())

	}

	rps := repository.NewRepository(db)
	srvces := service.NewService(rps)
	hdlrs := handler2.NewHandler(srvces)
	//handlers := new(handler2.Handler)
	srv := new(REST_JWT.Server)
	err = srv.Run(viper.GetString("port"), hdlrs.InitRoutes())
	if err != nil {
		log.Fatalf("ERROR on starting server: %s", err.Error())
	}
}
func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
