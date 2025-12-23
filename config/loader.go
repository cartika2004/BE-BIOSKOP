package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	JWT struct {
		Secret string
	}
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigName("config") 
	viper.SetConfigType("yaml")   
	viper.AddConfigPath(".")      

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}
	log.Println(" Config Loaded Successfully!")
}